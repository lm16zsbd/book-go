package user

import (
	"net/http"
	"strconv"

	"serica-go/internal/data"
	"serica-go/internal/server/middleware"
)

func GetUserID(r *http.Request) int64 {
	auth := middleware.GetAuth(r.Context())
	if auth != nil {
		return auth.UserID
	}
	return 0
}

func QueryInt(r *http.Request, key string, defaultVal int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return parsed
}

func QueryBool(r *http.Request, key string, defaultVal bool) bool {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	return v == "true" || v == "1"
}

func PathInt(r *http.Request, key string) int64 {
	v := r.PathValue(key)
	if v == "" {
		v = r.URL.Query().Get(key)
	}
	if v == "" {
		return 0
	}
	parsed, _ := strconv.ParseInt(v, 10, 64)
	return parsed
}

type PageResult struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	PageIndex  int         `json:"pageIndex"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

func NewPageResult(items interface{}, total int64, pageIndex, pageSize int) PageResult {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return PageResult{Items: items, Total: total, PageIndex: pageIndex, PageSize: pageSize, TotalPages: totalPages}
}

type BookDetailResponse struct {
	Book        *data.Book           `json:"book"`
	IsFavourite bool                 `json:"isFavourite"`
	Progress    *data.BookReadingPos `json:"progress,omitempty"`
	Bookmarks   []data.Bookmark      `json:"bookmarks,omitempty"`
	Annotations []data.Annotation    `json:"annotations,omitempty"`
	BookNotes   []data.BookNote      `json:"bookNotes,omitempty"`
}
