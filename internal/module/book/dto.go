package book

import (
	"encoding/json"
	"time"
)

type BookDetail struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	CoverUrl     string  `json:"coverUrl"`
	Url          string  `json:"url"`
	WrapKey      string  `json:"wrapKey"`
	Publisher    string  `json:"publisher"`
	PublishedAt  string  `json:"publishedAt"`
	ISBN         string  `json:"isbn"`
	FileType     string  `json:"fileType"`
	Language     string  `json:"language"`
	Desc         string  `json:"desc"`
	IsFavourite  bool    `json:"isFavourite"`
	Process      float64 `json:"process"`
	LastPosition string  `json:"lastPosition"`
}

func NewBookDetail(book *BookDetail, publishDate *time.Time) *BookDetail {
	if publishDate != nil {
		book.PublishedAt = publishDate.Format(time.RFC3339)
	}
	return book
}

type SelectionResponse struct {
	Language    []string `json:"language"`
	FileType    []string `json:"fileType"`
	Category    []string `json:"category"`
	PublishedAt []int    `json:"publishedAt"`
}

type CreateBookNoteReq struct {
	ID     *int64          `json:"id,omitempty"`
	BookID int64           `json:"bookId" validate:"required"`
	Note   json.RawMessage `json:"note" validate:"required"`
}

type UpdateBookNoteReq struct {
	Note json.RawMessage `json:"note" validate:"required"`
}

type TransReq struct {
	Data         string `json:"data" validate:"required"`
	SourceLang   string `json:"sourceLang"`
	TargetLang   string `json:"targetLang" validate:"required"`
	TargetScript string `json:"targetScript"`
}

type AuthorItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type PublisherItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}
