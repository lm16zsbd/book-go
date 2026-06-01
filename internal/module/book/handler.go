package book

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"serica-go/internal/conf"
	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/httputil"
)

type Handler struct {
	bookRepo *data.BookRepo
	userRepo *data.UserRepo
	cfg      *conf.Bootstrap
}

func NewHandler(bookRepo *data.BookRepo, userRepo *data.UserRepo, cfg *conf.Bootstrap) *Handler {
	return &Handler{bookRepo: bookRepo, userRepo: userRepo, cfg: cfg}
}

// @Summary      获取书籍列表
// @Description  分页获取书籍
// @Tags         Book
// @Produce      json
// @Param        keyword   query  string  false  "搜索关键词"
// @Param        pageIndex query  int     false  "页码"  default(1)
// @Param        pageSize  query  int     false  "每页数量"  default(20)
// @Param        order     query  string  false  "排序"
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/books [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageIndex := u.QueryInt(r, "pageIndex", 1)
	pageSize := u.QueryInt(r, "pageSize", 20)
	keyword := r.URL.Query().Get("keyword")
	order := r.URL.Query().Get("order")

	offset := (pageIndex - 1) * pageSize
	books, total, err := h.bookRepo.Paginate(data.BookFilter{
		Keyword: keyword,
		Order:   order,
	}, offset, pageSize)
	if err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}

	httputil.RespondJSON(w, 200, u.NewPageResult(books, total, pageIndex, pageSize))
}

// @Summary      获取书籍详情
// @Tags         Book
// @Produce      json
// @Param        id path int true "书籍ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /v1/client/books/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	bookID := u.PathInt(r, "id")
	userID := u.GetUserID(r)

	book, err := h.bookRepo.FindByID(bookID)
	if err != nil {
		httputil.RespondJSON(w, 404, map[string]string{"error": "book not found"})
		return
	}

	detail := u.BookDetailResponse{Book: book}

	if userID > 0 {
		fav, _ := h.userRepo.FindFavourite(userID, bookID)
		detail.IsFavourite = fav != nil

		pos, _ := h.bookRepo.GetReadingPos(userID, bookID)
		if pos != nil {
			detail.Progress = pos
		}

		bookmarks, _ := h.userRepo.FindBookmarks(userID, bookID)
		detail.Bookmarks = bookmarks

		annotations, _ := h.userRepo.FindAnnotations(userID, bookID)
		detail.Annotations = annotations

		notes, _ := h.userRepo.FindBookNotes(userID, bookID)
		detail.BookNotes = notes
	}

	httputil.RespondJSON(w, 200, detail)
}

// @Summary      搜索书籍
// @Tags         Book
// @Produce      json
// @Param        keyword   query  string  false  "搜索关键词"
// @Param        pageIndex query  int     false  "页码"  default(1)
// @Param        pageSize  query  int     false  "每页数量"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/books/search [get]
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	pageIndex := u.QueryInt(r, "pageIndex", 1)
	pageSize := u.QueryInt(r, "pageSize", 20)
	keyword := r.URL.Query().Get("keyword")
	category := r.URL.Query().Get("category")
	fileType := r.URL.Query().Get("fileType")
	publishYear := r.URL.Query().Get("publishYear")
	desc := !(r.URL.Query().Get("desc") == "false")

	offset := (pageIndex - 1) * pageSize
	order := "DESC"
	if !desc {
		order = "ASC"
	}
	books, total, err := h.bookRepo.Paginate(data.BookFilter{
		Keyword:  keyword,
		Category: parseCategory(category),
		Order:    order,
	}, offset, pageSize)
	if err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	_ = fileType
	_ = publishYear

	httputil.RespondJSON(w, 200, u.NewPageResult(books, total, pageIndex, pageSize))
}

func parseCategory(s string) int64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

// @Summary      获取分类列表
// @Tags         Book
// @Produce      json
// @Success      200  {array}   data.Category
// @Router       /v1/client/categories [get]
func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.bookRepo.ListCategories()
	if err != nil {
		httputil.RespondJSON(w, 200, []interface{}{})
		return
	}
	httputil.RespondJSON(w, 200, cats)
}

func (h *Handler) CategoriesAll(w http.ResponseWriter, r *http.Request) {
	cats, err := h.bookRepo.ListCategories()
	if err != nil {
		httputil.RespondJSON(w, 200, []interface{}{})
		return
	}
	httputil.RespondJSON(w, 200, cats)
}

func (h *Handler) CategoriesPaginated(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	pageIndex := u.QueryInt(r, "page", 1)
	pageSize := u.QueryInt(r, "limit", 20)
	offset := (pageIndex - 1) * pageSize

	cats, total, err := h.bookRepo.FindCategoriesPaginated(keyword, offset, pageSize)
	if err != nil {
		httputil.RespondJSON(w, 200, map[string]interface{}{"items": []interface{}{}, "total": 0, "page": pageIndex, "limit": pageSize})
		return
	}
	httputil.RespondJSON(w, 200, map[string]interface{}{
		"items": cats,
		"total": total,
		"page":  pageIndex,
		"limit": pageSize,
	})
}

func (h *Handler) CategoriesGetByID(w http.ResponseWriter, r *http.Request) {
	id := u.PathInt(r, "id")
	cat, err := h.bookRepo.FindCategoryByID(id)
	if err != nil {
		httputil.RespondJSON(w, 404, map[string]string{"error": "category not found"})
		return
	}
	httputil.RespondJSON(w, 200, cat)
}

func (h *Handler) CategoryBooks(w http.ResponseWriter, r *http.Request) {
	id := u.PathInt(r, "id")
	pageIndex := u.QueryInt(r, "page", 1)
	pageSize := u.QueryInt(r, "limit", 20)
	offset := (pageIndex - 1) * pageSize

	books, total, err := h.bookRepo.FindBooksByCategory(id, offset, pageSize)
	if err != nil {
		httputil.RespondJSON(w, 200, u.NewPageResult([]data.Book{}, 0, pageIndex, pageSize))
		return
	}
	httputil.RespondJSON(w, 200, u.NewPageResult(books, total, pageIndex, pageSize))
}

func (h *Handler) Authors(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	pageIndex := u.QueryInt(r, "page", 1)
	pageSize := u.QueryInt(r, "limit", 20)

	names, err := h.bookRepo.FindDistinctAuthors()
	if err != nil {
		httputil.RespondJSON(w, 200, map[string]interface{}{"items": []interface{}{}, "total": 0, "page": 1, "limit": 20})
		return
	}

	if name != "" {
		filtered := make([]string, 0)
		for _, n := range names {
			if len(n) >= len(name) && n[:len(name)] == name {
				filtered = append(filtered, n)
			}
		}
		names = filtered
	}

	total := int64(len(names))
	offset := (pageIndex - 1) * pageSize
	end := offset + pageSize
	if end > len(names) {
		end = len(names)
	}
	if offset >= len(names) {
		offset = 0
		end = 0
	}
	pagedNames := names[offset:end]

	items := make([]map[string]interface{}, len(pagedNames))
	for i, n := range pagedNames {
		count, _ := h.bookRepo.CountBooksByAuthor(n)
		items[i] = map[string]interface{}{"name": n, "count": count}
	}

	httputil.RespondJSON(w, 200, map[string]interface{}{
		"items": items,
		"total": total,
		"page":  pageIndex,
		"limit": pageSize,
	})
}

func (h *Handler) AuthorBooks(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	pageIndex := u.QueryInt(r, "page", 1)
	pageSize := u.QueryInt(r, "limit", 20)
	offset := (pageIndex - 1) * pageSize

	books, total, err := h.bookRepo.FindBooksByAuthor(name, offset, pageSize)
	if err != nil {
		httputil.RespondJSON(w, 200, u.NewPageResult([]data.Book{}, 0, pageIndex, pageSize))
		return
	}
	httputil.RespondJSON(w, 200, u.NewPageResult(books, total, pageIndex, pageSize))
}

func (h *Handler) Publishers(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	pageIndex := u.QueryInt(r, "page", 1)
	pageSize := u.QueryInt(r, "limit", 20)

	names, err := h.bookRepo.FindDistinctPublishers()
	if err != nil {
		httputil.RespondJSON(w, 200, map[string]interface{}{"items": []interface{}{}, "total": 0, "page": 1, "limit": 20})
		return
	}

	if name != "" {
		filtered := make([]string, 0)
		for _, n := range names {
			if len(n) >= len(name) && n[:len(name)] == name {
				filtered = append(filtered, n)
			}
		}
		names = filtered
	}

	total := int64(len(names))
	offset := (pageIndex - 1) * pageSize
	end := offset + pageSize
	if end > len(names) {
		end = len(names)
	}
	if offset >= len(names) {
		offset = 0
		end = 0
	}
	pagedNames := names[offset:end]

	items := make([]map[string]interface{}, len(pagedNames))
	for i, n := range pagedNames {
		count, _ := h.bookRepo.CountBooksByPublisher(n)
		items[i] = map[string]interface{}{"name": n, "count": count}
	}

	httputil.RespondJSON(w, 200, map[string]interface{}{
		"items": items,
		"total": total,
		"page":  pageIndex,
		"limit": pageSize,
	})
}

func (h *Handler) PublisherBooks(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	pageIndex := u.QueryInt(r, "page", 1)
	pageSize := u.QueryInt(r, "limit", 20)
	offset := (pageIndex - 1) * pageSize

	books, total, err := h.bookRepo.FindBooksByPublisher(name, offset, pageSize)
	if err != nil {
		httputil.RespondJSON(w, 200, u.NewPageResult([]data.Book{}, 0, pageIndex, pageSize))
		return
	}
	httputil.RespondJSON(w, 200, u.NewPageResult(books, total, pageIndex, pageSize))
}

// @Summary      获取筛选选项
// @Tags         Book
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/selection [get]
func (h *Handler) Selection(w http.ResponseWriter, r *http.Request) {
	cats, _ := h.bookRepo.ListCategories()
	years, _ := h.bookRepo.FindSelectionYears()

	categoryNames := make([]string, len(cats))
	for i, c := range cats {
		categoryNames[i] = c.Name
	}

	httputil.RespondJSON(w, 200, map[string]interface{}{
		"language":    []string{"繁體中文", "英文", "簡體中文"},
		"fileType":    []string{"pdf", "epub"},
		"category":    categoryNames,
		"publishedAt": years,
	})
}

func (h *Handler) GetReadingPos(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	bookID := u.PathInt(r, "bookId")

	pos, err := h.bookRepo.GetReadingPos(userID, bookID)
	if err != nil {
		httputil.RespondJSON(w, 200, map[string]interface{}{"process": 0, "lastPosition": ""})
		return
	}
	httputil.RespondJSON(w, 200, pos)
}

func (h *Handler) ReportReadingPos(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var pos data.BookReadingPos
	json.NewDecoder(r.Body).Decode(&pos)
	pos.UserID = userID
	if err := h.bookRepo.UpsertReadingPos(&pos); err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}

func (h *Handler) BookNotesList(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	bookID := int64(u.QueryInt(r, "bookId", 0))
	notes, _ := h.userRepo.FindBookNotes(userID, bookID)
	httputil.RespondJSON(w, 200, notes)
}

type CreateBookNoteInput struct {
	ID     *int64          `json:"id,omitempty"`
	BookID int64           `json:"bookId"`
	Note   json.RawMessage `json:"note"`
}

func (h *Handler) BookNotesCreate(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var input CreateBookNoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	exists, _ := h.bookRepo.ExistsByID(input.BookID)
	if !exists {
		httputil.RespondJSON(w, 404, map[string]interface{}{
			"code":    20001,
			"message": fmt.Sprintf("Book #%d not found", input.BookID),
		})
		return
	}
	note := &data.BookNote{
		UserID: userID,
		BookID: input.BookID,
		Note:   input.Note,
	}
	if input.ID != nil {
		note.ID = *input.ID
	}
	if err := h.userRepo.UpsertBookNote(note); err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, note)
}

type UpdateBookNoteInput struct {
	Note json.RawMessage `json:"note"`
}

func (h *Handler) BookNotesUpdate(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	id := u.PathInt(r, "id")
	var input UpdateBookNoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	note, err := h.userRepo.UpdateBookNote(id, userID, input.Note)
	if err != nil {
		httputil.RespondJSON(w, 404, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, note)
}

func (h *Handler) BookNotesDelete(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	id := u.PathInt(r, "id")
	if err := h.userRepo.DeleteBookNote(id, userID); err != nil {
		httputil.RespondJSON(w, 404, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, map[string]interface{}{"id": id})
}

func (h *Handler) BookNotesDeleteByPost(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	id := u.PathInt(r, "id")
	if err := h.userRepo.DeleteBookNote(id, userID); err != nil {
		httputil.RespondJSON(w, 404, map[string]string{"error": err.Error()})
		return
	}
	httputil.RespondJSON(w, 200, map[string]interface{}{"id": id})
}
