package home

import (
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/httputil"
	"serica-go/internal/utl"
)

type Handler struct {
	bookRepo *data.BookRepo
	userRepo *data.UserRepo
	redis    *data.RedisClient
}

func NewHandler(bookRepo *data.BookRepo, userRepo *data.UserRepo, redis *data.RedisClient) *Handler {
	return &Handler{bookRepo: bookRepo, userRepo: userRepo, redis: redis}
}

var homeBannerIds = []int64{103, 202, 183, 153, 193, 163}

// @Summary      首页数据
// @Description  获取首页 Banner 和分类
// @Tags         Home
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/client/home [get]
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	books, _ := h.bookRepo.FindByIDs(homeBannerIds)
	cats, _ := h.bookRepo.ListCategories()

	userID := u.GetUserID(r)
	bannerBooks := make([]map[string]interface{}, 0, len(books))
	for _, b := range books {
		isFav := false
		if userID > 0 {
			fav, err := h.userRepo.FindFavourite(userID, b.ID)
			isFav = err == nil && fav != nil
		}
		bannerBooks = append(bannerBooks, map[string]interface{}{
			"id":          b.ID,
			"title":       utl.FirstLang(b.Title),
			"author":      utl.FirstLang(b.Author),
			"coverUrl":    b.CoverUrl,
			"desc":        b.Desc,
			"isFavourite": isFav,
		})
	}

	httputil.RespondJSON(w, 200, map[string]interface{}{
		"banners":    bannerBooks,
		"categories": cats,
	})
}

// @Summary      首页栏目
// @Tags         Home
// @Produce      json
// @Success      200  {array}   map[string]interface{}
// @Router       /v1/client/homepage [get]
func (h *Handler) Homepage(w http.ResponseWriter, r *http.Request) {
	recommended, recTotal, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, 0, 10)
	newest, _, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, 0, 10)

	sections := []map[string]interface{}{
		{
			"title": "電子書推介",
			"items": recommended,
			"total": recTotal,
		},
		{
			"title": "最新上架",
			"items": newest,
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
