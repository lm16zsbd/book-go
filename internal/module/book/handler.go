package book

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"serica-go/internal/conf"
	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/binder"
	"serica-go/internal/pkg/dto"
	"serica-go/internal/pkg/exception"
	"serica-go/internal/pkg/httputil"
	"serica-go/internal/utl"
	"serica-go/internal/utl/crypto"
)

type Handler struct {
	bookRepo *data.BookRepo
	userRepo *data.UserRepo
	redis    *data.RedisClient
	cfg      *conf.Bootstrap
}

func NewHandler(bookRepo *data.BookRepo, userRepo *data.UserRepo, redis *data.RedisClient, cfg *conf.Bootstrap) *Handler {
	return &Handler{bookRepo: bookRepo, userRepo: userRepo, redis: redis, cfg: cfg}
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
	var q BookListQuery
	binder.BindQuery(r, &q)

	books, total, err := h.bookRepo.Paginate(q.ToFilter(), q.Offset(), q.PageSize)
	if err != nil {
		exception.InternalServerError.Write(w, err.Error())
		return
	}

	httputil.RespondJSON(w, 200, dto.NewPageResult(books, total, q.PageIndex, q.PageSize))
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
		BookNotFound.Write(w)
		return
	}

	result := BookDetail{
		ID:        book.ID,
		Title:     utl.FirstLang(book.Title),
		Author:    utl.FirstLang(book.Author),
		CoverUrl:  book.CoverUrl,
		Url:       book.Url,
		WrapKey:   h.getWrapKey(r.Context(), bookID, book.Key),
		Publisher: book.Publisher,
		ISBN:      book.ISBN,
		FileType:  book.FileType,
		Language:  book.Language,
		Desc:      book.Desc,
	}

	if book.PublishDate != nil {
		result.PublishedAt = book.PublishDate.Format(time.RFC3339)
	}

	if userID > 0 {
		if fav, err := h.userRepo.FindFavourite(userID, bookID); err == nil && fav != nil {
			result.IsFavourite = true
		}
		if pos, err := h.bookRepo.GetReadingPos(userID, bookID); err == nil && pos != nil {
			result.Process = pos.Process
			result.LastPosition = pos.LastPosition
		}
	}

	httputil.RespondJSON(w, 200, result)
}

func (h *Handler) getWrapKey(ctx context.Context, bookID int64, aesKey string) string {
	cacheKey := fmt.Sprintf("book:wrapkey:v3:%d", bookID)
	if h.redis != nil {
		if cached, err := h.redis.Get(ctx, cacheKey); err == nil && cached != "" {
			return cached
		}
	}

	if aesKey == "" {
		aesKey = os.Getenv("AES_KEY")
	}
	if aesKey == "" {
		return ""
	}

	rsaPubKey := os.Getenv("RSA_PUB_KEY")
	if rsaPubKey == "" {
		if h.redis != nil {
			_ = h.redis.Set(ctx, cacheKey, aesKey, 7*24*time.Hour)
		}
		return aesKey
	}

	wrap, err := crypto.RSAEncrypt(rsaPubKey, aesKey)
	if err != nil || wrap == "" {
		if h.redis != nil {
			_ = h.redis.Set(ctx, cacheKey, aesKey, 7*24*time.Hour)
		}
		return aesKey
	}

	if h.redis != nil {
		_ = h.redis.Set(ctx, cacheKey, wrap, 7*24*time.Hour)
	}
	return wrap
}

// @Summary      搜索书籍
// @Tags         Book
// @Produce      json
// @Param        keyword   query  string  false  "搜索关键词"
// @Param        pageIndex query  int     false  "页码"  default(1)
// @Param        pageSize  query  int     false  "每页数量"  default(20)
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/search [get]
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	var q BookListQuery
	binder.BindQuery(r, &q)

	books, total, err := h.bookRepo.Paginate(q.ToFilter(), q.Offset(), q.PageSize)
	if err != nil {
		exception.InternalServerError.Write(w, err.Error())
		return
	}

	httputil.RespondJSON(w, 200, dto.NewPageResult(books, total, q.PageIndex, q.PageSize))
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

// @Summary      获取所有分类
// @Tags         Book
// @Produce      json
// @Success      200  {array}   data.Category
// @Router       /v1/client/categories/all [get]
func (h *Handler) CategoriesAll(w http.ResponseWriter, r *http.Request) {
	cats, err := h.bookRepo.ListCategories()
	if err != nil {
		httputil.RespondJSON(w, 200, []interface{}{})
		return
	}
	httputil.RespondJSON(w, 200, cats)
}

func (h *Handler) CategoriesPaginated(w http.ResponseWriter, r *http.Request) {
	var q CategoriesPaginatedQuery
	binder.BindQuery(r, &q)

	cats, total, err := h.bookRepo.FindCategoriesPaginated(q.Keyword, q.Offset(), q.PageSize)
	if err != nil {
		httputil.RespondJSON(w, 200, dto.NewPageResult([]data.Category{}, 0, q.PageIndex, q.PageSize))
		return
	}
	httputil.RespondJSON(w, 200, dto.NewPageResult(cats, total, q.PageIndex, q.PageSize))
}

// @Summary      获取单个分类
// @Tags         Book
// @Produce      json
// @Param        id path int true "分类ID"
// @Success      200  {object}  data.Category
// @Router       /v1/client/categories/{id} [get]
func (h *Handler) CategoriesGetByID(w http.ResponseWriter, r *http.Request) {
	id := u.PathInt(r, "id")
	cat, err := h.bookRepo.FindCategoryByID(id)
	if err != nil {
		CategoryNotFound.Write(w)
		return
	}
	httputil.RespondJSON(w, 200, cat)
}

// @Summary      获取分类下书籍
// @Tags         Book
// @Produce      json
// @Param        id   path  int   true  "分类ID"
// @Param        page query int   false "页码" default(1)
// @Param        limit query int  false "每页数量" default(20)
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/categories/{id}/books [get]
func (h *Handler) CategoryBooks(w http.ResponseWriter, r *http.Request) {
	id := u.PathInt(r, "id")
	var q PageQuery
	binder.BindQuery(r, &q)

	books, total, err := h.bookRepo.FindBooksByCategory(id, q.Offset(), q.PageSize)
	if err != nil {
		httputil.RespondJSON(w, 200, dto.NewPageResult([]data.Book{}, 0, q.PageIndex, q.PageSize))
		return
	}
	httputil.RespondJSON(w, 200, dto.NewPageResult(books, total, q.PageIndex, q.PageSize))
}

// @Summary      获取作者列表
// @Tags         Book
// @Produce      json
// @Param        name  query string false "筛选名称"
// @Param        page  query int    false "页码" default(1)
// @Param        limit query int    false "每页数量" default(20)
// @Success      200   {object}  map[string]interface{}
// @Router       /v1/client/authors [get]
func (h *Handler) Authors(w http.ResponseWriter, r *http.Request) {
	var q AuthorsQuery
	binder.BindQuery(r, &q)

	names, err := h.bookRepo.FindDistinctAuthors()
	if err != nil {
		httputil.RespondJSON(w, 200, dto.NewPageResult([]AuthorItem{}, 0, q.PageIndex, q.PageSize))
		return
	}

	if q.Name != "" {
		filtered := make([]string, 0)
		for _, n := range names {
			if len(n) >= len(q.Name) && n[:len(q.Name)] == q.Name {
				filtered = append(filtered, n)
			}
		}
		names = filtered
	}

	total := int64(len(names))
	offset := q.Offset()
	end := offset + q.PageSize
	if end > len(names) {
		end = len(names)
	}
	if offset >= len(names) {
		offset = 0
		end = 0
	}
	pagedNames := names[offset:end]

	items := make([]AuthorItem, len(pagedNames))
	for i, n := range pagedNames {
		count, _ := h.bookRepo.CountBooksByAuthor(n)
		items[i] = AuthorItem{Name: n, Count: count}
	}

	httputil.RespondJSON(w, 200, dto.NewPageResult(items, total, q.PageIndex, q.PageSize))
}

// @Summary      获取作者下书籍
// @Tags         Book
// @Produce      json
// @Param        name path string true "作者名"
// @Param        page query int false "页码" default(1)
// @Param        limit query int false "每页数量" default(20)
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/authors/{name}/books [get]
func (h *Handler) AuthorBooks(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var q PageQuery
	binder.BindQuery(r, &q)

	books, total, err := h.bookRepo.FindBooksByAuthor(name, q.Offset(), q.PageSize)
	if err != nil {
		httputil.RespondJSON(w, 200, dto.NewPageResult([]data.Book{}, 0, q.PageIndex, q.PageSize))
		return
	}
	httputil.RespondJSON(w, 200, dto.NewPageResult(books, total, q.PageIndex, q.PageSize))
}

// @Summary      获取出版社列表
// @Tags         Book
// @Produce      json
// @Param        name  query string false "筛选名称"
// @Param        page  query int    false "页码" default(1)
// @Param        limit query int    false "每页数量" default(20)
// @Success      200   {object}  map[string]interface{}
// @Router       /v1/client/publishers [get]
func (h *Handler) Publishers(w http.ResponseWriter, r *http.Request) {
	var q AuthorsQuery
	binder.BindQuery(r, &q)

	names, err := h.bookRepo.FindDistinctPublishers()
	if err != nil {
		httputil.RespondJSON(w, 200, dto.NewPageResult([]PublisherItem{}, 0, q.PageIndex, q.PageSize))
		return
	}

	if q.Name != "" {
		filtered := make([]string, 0)
		for _, n := range names {
			if len(n) >= len(q.Name) && n[:len(q.Name)] == q.Name {
				filtered = append(filtered, n)
			}
		}
		names = filtered
	}

	total := int64(len(names))
	offset := q.Offset()
	end := offset + q.PageSize
	if end > len(names) {
		end = len(names)
	}
	if offset >= len(names) {
		offset = 0
		end = 0
	}
	pagedNames := names[offset:end]

	items := make([]PublisherItem, len(pagedNames))
	for i, n := range pagedNames {
		count, _ := h.bookRepo.CountBooksByPublisher(n)
		items[i] = PublisherItem{Name: n, Count: count}
	}

	httputil.RespondJSON(w, 200, dto.NewPageResult(items, total, q.PageIndex, q.PageSize))
}

// @Summary      获取出版社下书籍
// @Tags         Book
// @Produce      json
// @Param        name path string true "出版社名"
// @Param        page query int false "页码" default(1)
// @Param        limit query int false "每页数量" default(20)
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/publishers/{name}/books [get]
func (h *Handler) PublisherBooks(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var q PageQuery
	binder.BindQuery(r, &q)

	books, total, err := h.bookRepo.FindBooksByPublisher(name, q.Offset(), q.PageSize)
	if err != nil {
		httputil.RespondJSON(w, 200, dto.NewPageResult([]data.Book{}, 0, q.PageIndex, q.PageSize))
		return
	}

	httputil.RespondJSON(w, 200, dto.NewPageResult(books, total, q.PageIndex, q.PageSize))
}

// @Summary      获取筛选选项
// @Tags         Book
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/selection [get]
func (h *Handler) Selection(w http.ResponseWriter, r *http.Request) {
	cacheKey := "books:selection"
	if h.redis != nil {
		var cached map[string]interface{}
		if err := h.redis.GetJSON(r.Context(), cacheKey, &cached); err == nil && cached != nil {
			httputil.RespondJSON(w, 200, cached)
			return
		}
	}

	cats, _ := h.bookRepo.FindDistinctCategories()
	years, _ := h.bookRepo.FindSelectionYears()

	result := SelectionResponse{
		Language:    []string{"繁體中文", "英文", "簡體中文"},
		FileType:    []string{"pdf", "epub"},
		Category:    cats,
		PublishedAt: years,
	}

	if h.redis != nil {
		_ = h.redis.SetJSON(r.Context(), cacheKey, result, time.Hour)
	}

	httputil.RespondJSON(w, 200, result)
}

// @Summary      获取阅读进度
// @Tags         Book
// @Produce      json
// @Param        bookId query int true "书籍ID"
// @Success      200  {object}  data.BookReadingPos
// @Security     BearerAuth
// @Router       /v1/client/reading-pos [get]
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

// @Summary      上报阅读进度
// @Tags         Book
// @Accept       json
// @Produce      json
// @Param        body body data.BookReadingPos true "阅读进度"
// @Success      200  {object}  map[string]bool
// @Security     BearerAuth
// @Router       /v1/client/reading-pos [post]
func (h *Handler) ReportReadingPos(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var pos data.BookReadingPos
	json.NewDecoder(r.Body).Decode(&pos)
	pos.UserID = userID
	if err := h.bookRepo.UpsertReadingPos(&pos); err != nil {
		exception.InternalServerError.Write(w, err.Error())
		return
	}
	httputil.RespondJSON(w, 200, map[string]bool{"ok": true})
}

// @Summary      获取书籍笔记列表
// @Tags         Book Notes
// @Produce      json
// @Param        bookId query int true "书籍ID"
// @Success      200  {array}   map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/book-notes [get]
func (h *Handler) BookNotesList(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var q BookIDQuery
	binder.BindQuery(r, &q)
	notes, _ := h.userRepo.FindBookNotes(userID, q.BookID)
	httputil.RespondJSON(w, 200, notes)
}

// @Summary      创建/更新书籍笔记
// @Tags         Book Notes
// @Accept       json
// @Produce      json
// @Param        body body CreateBookNoteReq true "笔记内容"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/book-notes [post]
func (h *Handler) BookNotesCreate(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	var input CreateBookNoteReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		exception.InvalidBody.Write(w)
		return
	}
	exists, _ := h.bookRepo.ExistsByID(input.BookID)
	if !exists {
		BookNotFound.Write(w, fmt.Sprintf("Book #%d not found", input.BookID))
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
		exception.InternalServerError.Write(w, err.Error())
		return
	}
	httputil.RespondJSON(w, 200, note)
}

// @Summary      更新书籍笔记
// @Tags         Book Notes
// @Accept       json
// @Produce      json
// @Param        id   path int                true "笔记ID"
// @Param        body body UpdateBookNoteReq true "笔记内容"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /v1/client/book-notes/{id} [patch]
func (h *Handler) BookNotesUpdate(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	id := u.PathInt(r, "id")
	var input UpdateBookNoteReq
	if err := httputil.DecodeAndValidate(r, &input); err != nil {
		exception.InvalidBody.Write(w)
		return
	}
	note, err := h.userRepo.UpdateBookNote(id, userID, input.Note)
	if err != nil {
		BookNotFound.Write(w, fmt.Sprintf("Book #%d not found", id))
		return
	}
	httputil.RespondJSON(w, 200, note)
}

// @Summary      删除书籍笔记
// @Tags         Book Notes
// @Produce      json
// @Param        id path int true "笔记ID"
// @Success      200  {object}  map[string]bool
// @Security     BearerAuth
// @Router       /v1/client/book-notes/{id} [delete]
func (h *Handler) BookNotesDelete(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	id := u.PathInt(r, "id")
	if err := h.userRepo.DeleteBookNote(id, userID); err != nil {
		BookNotFound.Write(w, fmt.Sprintf("Book #%d not found", id))
		return
	}
	httputil.RespondJSON(w, 200, map[string]interface{}{"id": id})
}

func (h *Handler) BookNotesDeleteByPost(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)
	id := u.PathInt(r, "id")
	if err := h.userRepo.DeleteBookNote(id, userID); err != nil {
		BookNotFound.Write(w, fmt.Sprintf("Book #%d not found", id))
		return
	}
	httputil.RespondJSON(w, 200, map[string]interface{}{"id": id})
}
