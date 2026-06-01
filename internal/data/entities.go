package data

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email       string    `gorm:"uniqueIndex" json:"email,omitempty"`
	PhoneNumber string    `gorm:"column:phone_number" json:"phoneNumber,omitempty"`
	Name        string    `json:"name,omitempty"`
	Password    string    `json:"-"`
	IsActive    bool      `gorm:"column:is_active;default:true" json:"isActive"`
	PatronID    string    `gorm:"column:patronid" json:"patronid"`
	LoginDate   time.Time `gorm:"column:login_date" json:"loginDate"`
	Avatar      string    `json:"avatar,omitempty"`
	Language    string    `gorm:"default:en" json:"language,omitempty"`
	Offset      string    `gorm:"default:+00" json:"offset,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string { return "users" }

type UserSession struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;index" json:"userId"`
	Token     string    `gorm:"uniqueIndex" json:"token"`
	IP        string    `gorm:"column:ip" json:"ip"`
	UserAgent string    `gorm:"column:user_agent" json:"userAgent"`
	Device    string    `json:"device"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (UserSession) TableName() string { return "user_sessions" }

type Book struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string     `json:"title"`
	Excerpt         string     `json:"excerpt"`
	StockQuantity   int        `gorm:"column:stock_quantity" json:"stockQuantity"`
	PublishingGroup string     `gorm:"column:publishing_group" json:"publishingGroup,omitempty"`
	Imprints        string     `json:"imprints,omitempty"`
	Author          string     `json:"author"`
	TotalPages      int        `gorm:"column:total_pages" json:"totalPages,omitempty"`
	Cost            string     `json:"cost,omitempty"`
	CoverPhoto      string     `gorm:"column:cover_photo" json:"coverPhoto"`
	BookPDF         string     `gorm:"column:book_pdf" json:"bookPdf"`
	FileType        string     `gorm:"column:file_type;default:pdf" json:"fileType,omitempty"`
	IsRecomm        bool       `gorm:"column:is_recomm" json:"isRecomm,omitempty"`
	PreviewBook     string     `gorm:"column:preview_book" json:"previewBook,omitempty"`
	AddedBy         int64      `gorm:"column:added_by" json:"addedBy"`
	AddedAt         time.Time  `gorm:"column:added_at" json:"addedAt"`
	PublishDate     *time.Time `gorm:"column:publish_date" json:"publishDate,omitempty"`
	ISBN            string     `gorm:"column:isbn_no" json:"isbnNo"`
	IsDeleted       bool       `gorm:"column:is_deleted;default:false" json:"isDeleted"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deletedAt,omitempty"`
	CollectionType  string     `gorm:"column:collection_type" json:"collectionType,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Book) TableName() string { return "books" }

type Favourite struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"column:user_id;index" json:"userId"`
	BookID    int64      `gorm:"column:book_id;index" json:"bookId"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deletedAt,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
}

func (Favourite) TableName() string { return "favourites" }

type Bookmark struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"userId"`
	BookID    int64     `gorm:"column:book_id" json:"bookId"`
	Title     string    `json:"title"`
	Position  string    `json:"position"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Bookmark) TableName() string { return "bookmarks" }

type Annotation struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"userId"`
	BookID    int64     `gorm:"column:book_id" json:"bookId"`
	Content   string    `json:"content"`
	Position  string    `json:"position"`
	Color     string    `json:"color"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Annotation) TableName() string { return "annotations" }

type BookNote struct {
	ID        int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64           `gorm:"column:user_id" json:"userId"`
	BookID    int64           `gorm:"column:book_id" json:"bookId"`
	Note      json.RawMessage `gorm:"type:jsonb" json:"note"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"deletedAt,omitempty"`
}

func (BookNote) TableName() string { return "book_notes" }

type BookReadingPos struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"column:user_id;uniqueIndex:uq_user_book" json:"userId"`
	BookID       int64     `gorm:"column:book_id;uniqueIndex:uq_user_book" json:"bookId"`
	Process      float64   `json:"process"`
	LastPosition string    `gorm:"column:last_position" json:"lastPosition"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (BookReadingPos) TableName() string { return "book_reading_pos" }

type ReaderConfig struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;uniqueIndex" json:"userId"`
	Theme     string    `gorm:"default:light" json:"theme"`
	FontSize  int       `gorm:"column:font_size;default:16" json:"fontSize"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ReaderConfig) TableName() string { return "reader_configs" }

type BookKey struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BookID    int64     `gorm:"column:book_id;uniqueIndex" json:"bookId"`
	Key       string    `json:"key"`
	URL       string    `json:"url,omitempty"`
	Status    string    `gorm:"default:PENDING" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (BookKey) TableName() string { return "book_keys" }

type Category struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Category) TableName() string { return "categories" }

type Banner struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ImageURL  string    `gorm:"column:image_url" json:"imageUrl"`
	BookID    int64     `gorm:"column:book_id" json:"bookId,omitempty"`
	Order     int       `json:"order"`
	IsActive  bool      `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Banner) TableName() string { return "banners" }

type SampleBook struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BookID    int64     `gorm:"column:book_id" json:"bookId"`
	URL       string    `json:"url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (SampleBook) TableName() string { return "sample_books" }

type SearchHistory struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"userId"`
	Keyword   string    `json:"keyword"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (SearchHistory) TableName() string { return "search_histories" }
