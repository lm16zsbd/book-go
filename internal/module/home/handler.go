package home

import (
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/binder"
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
// @Success      200  {object}  home.HomeResponse
// @Router       /v1/client/home [get]
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	books, _ := h.bookRepo.FindByIDs(homeBannerIds)
	cats, _ := h.bookRepo.ListCategories()

	userID := u.GetUserID(r)
	bannerBooks := make([]BannerItem, 0, len(books))
	for _, b := range books {
		isFav := false
		if userID > 0 {
			fav, err := h.userRepo.FindFavourite(userID, b.ID)
			isFav = err == nil && fav != nil
		}
		bannerBooks = append(bannerBooks, BannerItem{
			ID:          b.ID,
			Title:       utl.FirstLang(b.Title),
			Author:      utl.FirstLang(b.Author),
			CoverUrl:    b.CoverUrl,
			Desc:        b.Desc,
			IsFavourite: isFav,
		})
	}

	httputil.RespondJSON(w, 200, HomeResponse{Banners: bannerBooks, Categories: cats})
}

// @Summary      首页栏目
// @Tags         Home
// @Produce      json
// @Success      200  {array}   home.HomepageSection
// @Router       /v1/client/homepage [get]
func (h *Handler) Homepage(w http.ResponseWriter, r *http.Request) {
	recommended, recTotal, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, 0, 10)
	newest, _, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, 0, 10)

	sections := []HomepageSection{
		{Title: "電子書推介", Items: recommended, Total: recTotal},
		{Title: "最新上架", Items: newest},
	}

	httputil.RespondJSON(w, 200, sections)
}

func (h *Handler) HomepageSection(w http.ResponseWriter, r *http.Request) {
	var q HomepageSectionQuery
	binder.BindQuery(r, &q)

	offset := (q.PageIndex - 1) * q.PageSize
	books, total, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, offset, q.PageSize)

	httputil.RespondJSON(w, 200, HomepageSectionResponse{
		Title:     q.Title,
		Items:     books,
		Total:     total,
		PageIndex: q.PageIndex,
		PageSize:  q.PageSize,
	})
}
