package data

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

func newRepos(db *gorm.DB) Repos {
	return Repos{
		User:          NewUserRepo(db),
		Book:          NewBookRepo(db),
		BookKey:       NewBookKeyRepo(db),
		SearchHistory: NewSearchHistoryRepo(db),
	}
}

// ----- User Repo -----
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByEmail(email string) (*User, error) {
	var u User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) CreateUser(u *User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) FindByID(id int64) (*User, error) {
	var u User
	err := r.db.First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepo) FindSessionByToken(token string) (*UserSession, error) {
	var s UserSession
	err := r.db.Where("token = ?", token).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *UserRepo) CreateSession(s *UserSession) error {
	return r.db.Create(s).Error
}

func (r *UserRepo) FindFavourite(userID, bookID int64) (*Favourite, error) {
	var f Favourite
	err := r.db.Where("user_id = ? AND book_id = ? AND deleted_at IS NULL", userID, bookID).First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *UserRepo) FindFavouritesPaginated(userID int64, offset, limit int, desc bool) ([]map[string]interface{}, int64, error) {
	var favourites []Favourite
	var total int64

	order := "DESC"
	if !desc {
		order = "ASC"
	}

	q := r.db.Model(&Favourite{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	q.Count(&total)
	q.Order("created_at " + order).Offset(offset).Limit(limit).Find(&favourites)

	bookIDs := make([]int64, len(favourites))
	for i, f := range favourites {
		bookIDs[i] = f.BookID
	}

	type favBook struct {
		ID           int64   `json:"id"`
		Title        string  `json:"title"`
		Author       string  `json:"author"`
		CoverPhoto   string  `json:"coverUrl"`
		Excerpt      string  `json:"desc"`
		Process      float64 `json:"process"`
		LastPosition string  `json:"lastPosition"`
	}

	var books []favBook
	r.db.Model(&Book{}).
		Select("books.id, books.title, books.author, books.cover_url, books.desc, COALESCE(bp.process, 0) as process, COALESCE(bp.last_position, '') as last_position").
		Joins("LEFT JOIN book_reading_pos bp ON bp.user_id = ? AND bp.book_id = books.id", userID).
		Where("books.id IN ?", bookIDs).
		Find(&books)

	bookMap := make(map[int64]favBook, len(books))
	for _, b := range books {
		bookMap[b.ID] = b
	}

	items := make([]map[string]interface{}, 0, len(favourites))
	for _, f := range favourites {
		if b, ok := bookMap[f.BookID]; ok {
			items = append(items, map[string]interface{}{
				"id":           b.ID,
				"title":        b.Title,
				"author":       b.Author,
				"coverUrl":     b.CoverPhoto,
				"desc":         b.Excerpt,
				"process":      b.Process,
				"lastPosition": b.LastPosition,
			})
		}
	}

	return items, total, nil
}

func (r *UserRepo) ToggleFavourite(userID, bookID int64) (bool, error) {
	fav, err := r.FindFavourite(userID, bookID)
	if err == nil && fav != nil {
		r.db.Model(fav).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))
		return false, nil
	}
	return true, r.db.Create(&Favourite{UserID: userID, BookID: bookID}).Error
}

func (r *UserRepo) BatchCancelFavourite(userID int64, bookIDs []int64) int64 {
	result := r.db.Model(&Favourite{}).
		Where("user_id = ? AND book_id IN ? AND deleted_at IS NULL", userID, bookIDs).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))
	return result.RowsAffected
}

func (r *UserRepo) BatchCancelAllFavourite(userID int64, excludeIDs []int64) int64 {
	q := r.db.Model(&Favourite{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)
	if len(excludeIDs) > 0 {
		q = q.Where("book_id NOT IN ?", excludeIDs)
	}
	result := q.Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))
	return result.RowsAffected
}

func (r *UserRepo) FindBookmarks(userID, bookID int64) ([]Bookmark, error) {
	var bookmarks []Bookmark
	err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Find(&bookmarks).Error
	return bookmarks, err
}

func (r *UserRepo) CreateBookmark(b *Bookmark) error {
	return r.db.Create(b).Error
}

func (r *UserRepo) DeleteBookmark(id int64) error {
	return r.db.Delete(&Bookmark{}, id).Error
}

func (r *UserRepo) FindAnnotations(userID, bookID int64) ([]Annotation, error) {
	var annotations []Annotation
	err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Find(&annotations).Error
	return annotations, err
}

func (r *UserRepo) CreateAnnotation(a *Annotation) error {
	return r.db.Create(a).Error
}

func (r *UserRepo) UpdateAnnotation(id int64, updates map[string]interface{}) error {
	return r.db.Model(&Annotation{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepo) DeleteAnnotation(id int64) error {
	return r.db.Delete(&Annotation{}, id).Error
}

func (r *UserRepo) GetReaderConfig(userID int64) (*ReaderConfig, error) {
	var cfg ReaderConfig
	err := r.db.Where("user_id = ?", userID).First(&cfg).Error
	return &cfg, err
}

func (r *UserRepo) SaveReaderConfig(cfg *ReaderConfig) error {
	return r.db.Save(cfg).Error
}

// ----- Book Repo -----
type BookRepo struct {
	db *gorm.DB
}

func NewBookRepo(db *gorm.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (r *BookRepo) FindByID(id int64) (*Book, error) {
	var b Book
	err := r.db.Where("id = ?", id).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookRepo) ExistsByID(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&Book{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

type BookFilter struct {
	Keyword      string
	Category     string
	CategoryIDs  []int64
	FileType     string
	FileTypes    []string
	PublishYear  string
	PublishYears []string
	Language     string
	Languages    []string
	Order        string
	Desc         bool
}

func (r *BookRepo) Paginate(filter BookFilter, offset, limit int) ([]Book, int64, error) {
	q := r.db.Model(&Book{})

	if filter.Keyword != "" {
		q = q.Where("title ILIKE ? OR author ILIKE ? OR isbn ILIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}
	if filter.Category != "" {
		q = q.Where("category ILIKE ?", "%"+filter.Category+"%")
	}
	if len(filter.FileTypes) > 0 {
		q = q.Where("file_type IN ?", filter.FileTypes)
	}
	if len(filter.PublishYears) > 0 {
		q = q.Where("EXTRACT(YEAR FROM published_at) IN ?", filter.PublishYears)
	}
	if len(filter.Languages) > 0 {
		q = q.Where("language IN ?", filter.Languages)
	}

	var total int64
	q.Count(&total)

	books := make([]Book, 0)
	order := "created_at DESC"
	if !filter.Desc {
		order = "created_at ASC"
	}
	if filter.Order == "title" {
		order = "title ASC"
	}
	q.Order(order).Offset(offset).Limit(limit).Find(&books)
	return books, total, nil
}

func (r *BookRepo) Search(keyword string, offset, limit int) ([]Book, int64, error) {
	q := r.db.Model(&Book{})
	if keyword != "" {
		q = q.Where("title ILIKE ? OR author ILIKE ? OR isbn ILIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	var total int64
	q.Count(&total)
	books := make([]Book, 0)
	q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&books)
	return books, total, nil
}

// ----- BookKey Repo -----
type BookKeyRepo struct {
	db *gorm.DB
}

func NewBookKeyRepo(db *gorm.DB) *BookKeyRepo {
	return &BookKeyRepo{db: db}
}

func (r *BookKeyRepo) FindByBookID(bookID int64) (*BookKey, error) {
	var bk BookKey
	err := r.db.Where("book_id = ?", bookID).First(&bk).Error
	if err != nil {
		return nil, err
	}
	return &bk, nil
}

func (r *BookKeyRepo) Upsert(bk *BookKey) error {
	return r.db.Save(bk).Error
}

// ----- ReadingPos Repo -----
func (r *BookRepo) GetReadingPos(userID, bookID int64) (*BookReadingPos, error) {
	var pos BookReadingPos
	err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&pos).Error
	if err != nil {
		return nil, err
	}
	return &pos, nil
}

func (r *BookRepo) UpsertReadingPos(pos *BookReadingPos) error {
	return r.db.Where("user_id = ? AND book_id = ?", pos.UserID, pos.BookID).
		Assign(map[string]interface{}{
			"process":       pos.Process,
			"last_position": pos.LastPosition,
		}).
		FirstOrCreate(pos).Error
}

// ----- Categories -----
func (r *BookRepo) ListCategories() ([]Category, error) {
	var cats []Category
	err := r.db.Find(&cats).Error
	return cats, err
}

// ----- Banners -----
func (r *BookRepo) ListBanners() ([]Banner, error) {
	var banners []Banner
	err := r.db.Where("is_active = ?", true).Order("\"order\" ASC").Find(&banners).Error
	return banners, err
}

func (r *BookRepo) FindDistinctAuthors() ([]string, error) {
	var names []string
	err := r.db.Model(&Book{}).Where("author != ''").Distinct("author").Pluck("author", &names).Error
	return names, err
}

func (r *BookRepo) FindDistinctPublishers() ([]string, error) {
	var names []string
	err := r.db.Model(&Book{}).Where("publisher != ''").Distinct("publisher").Pluck("publisher", &names).Error
	return names, err
}

func (r *BookRepo) CountBooksByAuthor(author string) (int64, error) {
	var count int64
	err := r.db.Model(&Book{}).Where("author = ?", author).Count(&count).Error
	return count, err
}

func (r *BookRepo) CountBooksByPublisher(publisher string) (int64, error) {
	var count int64
	err := r.db.Model(&Book{}).Where("publisher = ?", publisher).Count(&count).Error
	return count, err
}

func (r *BookRepo) FindBooksByAuthor(author string, offset, limit int) ([]Book, int64, error) {
	var total int64
	q := r.db.Model(&Book{}).Where("author = ?", author)
	q.Count(&total)
	books := make([]Book, 0)
	q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&books)
	return books, total, nil
}

func (r *BookRepo) FindBooksByPublisher(publisher string, offset, limit int) ([]Book, int64, error) {
	var total int64
	q := r.db.Model(&Book{}).Where("publisher = ?", publisher)
	q.Count(&total)
	books := make([]Book, 0)
	q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&books)
	return books, total, nil
}

func (r *BookRepo) FindCategoriesPaginated(keyword string, offset, limit int) ([]Category, int64, error) {
	var total int64
	q := r.db.Model(&Category{})
	if keyword != "" {
		q = q.Where("name ILIKE ?", "%"+keyword+"%")
	}
	q.Count(&total)
	var cats []Category
	q.Order("created_at ASC").Offset(offset).Limit(limit).Find(&cats)
	return cats, total, nil
}

func (r *BookRepo) FindCategoryByID(id int64) (*Category, error) {
	var cat Category
	err := r.db.First(&cat, id).Error
	return &cat, err
}

func (r *BookRepo) FindBooksByCategory(categoryID int64, offset, limit int) ([]Book, int64, error) {
	var total int64
	q := r.db.Model(&Book{})
	q = q.Where("id IN (SELECT book_id FROM book_categories WHERE category_id = ?)", categoryID)
	q.Count(&total)
	books := make([]Book, 0)
	q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&books)
	return books, total, nil
}

func (r *BookRepo) FindSelectionYears() ([]int, error) {
	var years []int
	err := r.db.Model(&Book{}).Where("published_at IS NOT NULL").
		Select("DISTINCT EXTRACT(YEAR FROM published_at)").
		Order("EXTRACT(YEAR FROM published_at) DESC").
		Pluck("EXTRACT(YEAR FROM published_at)", &years).Error
	return years, err
}

// ----- Book Notes -----
func (r *UserRepo) FindBookNotes(userID, bookID int64) ([]BookNote, error) {
	var notes []BookNote
	err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Order("created_at DESC").Find(&notes).Error
	return notes, err
}

func (r *UserRepo) CreateBookNote(note *BookNote) error {
	return r.db.Create(note).Error
}

func (r *UserRepo) UpsertBookNote(note *BookNote) error {
	if note.ID > 0 {
		return r.db.Model(&BookNote{}).Where("id = ? AND user_id = ?", note.ID, note.UserID).Updates(map[string]interface{}{
			"note": note.Note,
		}).Error
	}
	return r.db.Create(note).Error
}

func (r *UserRepo) UpdateBookNote(id, userID int64, note json.RawMessage) (*BookNote, error) {
	var existing BookNote
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&existing).Error; err != nil {
		return nil, fmt.Errorf("book note not found")
	}
	existing.Note = note
	if err := r.db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *UserRepo) DeleteBookNote(id, userID int64) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&BookNote{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("book note not found")
	}
	return nil
}

// ----- Search History Repo -----
type SearchHistoryRepo struct {
	db *gorm.DB
}

func NewSearchHistoryRepo(db *gorm.DB) *SearchHistoryRepo {
	return &SearchHistoryRepo{db: db}
}

func (r *SearchHistoryRepo) FindByUser(userID int64, limit int) ([]SearchHistory, error) {
	var histories []SearchHistory
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&histories).Error
	return histories, err
}

func (r *SearchHistoryRepo) Create(history *SearchHistory) error {
	return r.db.Create(history).Error
}

func (r *SearchHistoryRepo) DeleteAll(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&SearchHistory{}).Error
}

func (r *SearchHistoryRepo) FindHotKeywords(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Model(&SearchHistory{}).
		Select("keyword, COUNT(*) as count").
		Group("keyword").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}
