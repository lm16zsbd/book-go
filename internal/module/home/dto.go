package home

import "serica-go/internal/data"

type BannerItem struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	CoverUrl    string `json:"coverUrl"`
	Desc        string `json:"desc"`
	IsFavourite bool   `json:"isFavourite"`
}

type HomeResponse struct {
	Banners    []BannerItem    `json:"banners"`
	Categories []data.Category `json:"categories"`
}

type HomepageSection struct {
	Title string      `json:"title"`
	Items []data.Book `json:"items"`
	Total int64       `json:"total,omitempty"`
}

type HomepageSectionResponse struct {
	Title     string      `json:"title"`
	Items     []data.Book `json:"items"`
	Total     int64       `json:"total"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}

type HomepageSectionQuery struct {
	Title     string `query:"title"`
	PageIndex int    `query:"pageIndex,1"`
	PageSize  int    `query:"pageSize,12"`
}
