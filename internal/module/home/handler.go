package home

import (
	"net/http"

	"serica-go/internal/data"
	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/binder"
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

var homeSectionIDs = [][]int64{
	{103, 202, 183, 153, 193, 163, 185, 200, 201, 199, 158, 186},
	{191, 131, 171, 179, 190, 157, 189, 143, 167, 188, 154, 180},
	{196, 197, 195, 192, 162, 194, 173, 164, 181, 136, 128, 198},
}

var homeSectionTitles = []string{
	"電子書推介",
	"最受歡迎電子書館藏",
	"自訂書籍展示模組",
}

// @Summary      首页数据
// @Description  获取首页 Banner 和分类
// @Tags         Home
// @Produce      json
// @Success      200  {object}  home.HomeResponse
// @Router       /v1/client/home [get]
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) (any, error) {
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

	return HomeResponse{Banners: bannerBooks, Categories: cats}, nil
}

// @Summary      首页栏目
// @Tags         Home
// @Produce      json
// @Success      200  {array}   home.HomepageSection
// @Router       /v1/client/homepage [get]
func (h *Handler) Homepage(w http.ResponseWriter, r *http.Request) (any, error) {
	sections := make([]HomepageSection, 0, len(homeSectionIDs))

	for i, ids := range homeSectionIDs {
		books, _ := h.bookRepo.FindByIDs(ids)

		title := ""
		if i < len(homeSectionTitles) {
			title = homeSectionTitles[i]
		}

		// Maintain original order from IDs
		ordered := make([]data.Book, 0, len(ids))
		bookMap := make(map[int64]data.Book, len(ids))
		for _, b := range books {
			bookMap[b.ID] = b
		}
		for _, id := range ids {
			if b, ok := bookMap[id]; ok {
				ordered = append(ordered, b)
			}
		}

		sections = append(sections, HomepageSection{
			Title: title,
			Items: ordered,
			Total: int64(len(ids)),
		})
	}

	return sections, nil
}

func (h *Handler) HomepageSection(w http.ResponseWriter, r *http.Request) (any, error) {
	var q HomepageSectionQuery
	binder.BindQuery(r, &q)

	offset := (q.PageIndex - 1) * q.PageSize
	books, total, _ := h.bookRepo.Paginate(data.BookFilter{Order: "DESC"}, offset, q.PageSize)

	return HomepageSectionResponse{
		Title:     q.Title,
		Items:     books,
		Total:     total,
		PageIndex: q.PageIndex,
		PageSize:  q.PageSize,
	}, nil
}
