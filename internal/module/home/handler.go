package home

import (
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/httputil"
)

type Handler struct {
	bookRepo *data.BookRepo
	userRepo *data.UserRepo
	redis    *data.RedisClient
}

func NewHandler(bookRepo *data.BookRepo, userRepo *data.UserRepo, redis *data.RedisClient) *Handler {
	return &Handler{bookRepo: bookRepo, userRepo: userRepo, redis: redis}
}

// @Summary      首页数据
// @Description  获取首页 Banner、分类和书籍列表
// @Tags         Home
// @Produce      json
// @Param        pageIndex query  int  false  "页码"  default(1)
// @Param        pageSize  query  int  false  "每页数量"  default(10)
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/home [get]
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	pageIndex := u.QueryInt(r, "pageIndex", 1)
	pageSize := u.QueryInt(r, "pageSize", 10)

	offset := (pageIndex - 1) * pageSize
	books, total, _ := h.bookRepo.Paginate(data.BookFilter{}, offset, pageSize)
	banners, _ := h.bookRepo.ListBanners()
	cats, _ := h.bookRepo.ListCategories()

	httputil.RespondJSON(w, 200, map[string]interface{}{
		"banners":    banners,
		"categories": cats,
		"books":      u.NewPageResult(books, total, pageIndex, pageSize),
	})
}

// @Summary      首页栏目
// @Tags         Home
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Router       /v1/client/homepage [get]
func (h *Handler) Homepage(w http.ResponseWriter, r *http.Request) {
	pageIndex := u.QueryInt(r, "pageIndex", 1)
	pageSize := u.QueryInt(r, "pageSize", 10)

	offset := (pageIndex - 1) * pageSize

	recommended, recTotal, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, offset, pageSize)
	newest, newTotal, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, 0, pageSize)

	sections := []map[string]interface{}{
		{
			"title": "電子書推介",
			"items": u.NewPageResult(recommended, recTotal, pageIndex, pageSize),
		},
		{
			"title": "最新上架",
			"items": u.NewPageResult(newest, newTotal, 1, pageSize),
		},
	}

	httputil.RespondJSON(w, 200, sections)
}

func (h *Handler) HomepageSection(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	pageIndex := u.QueryInt(r, "pageIndex", 1)
	pageSize := u.QueryInt(r, "pageSize", 12)

	offset := (pageIndex - 1) * pageSize
	books, total, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, offset, pageSize)

	httputil.RespondJSON(w, 200, map[string]interface{}{
		"title":     title,
		"items":     books,
		"total":     total,
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	})
}
