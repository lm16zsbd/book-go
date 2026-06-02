package book

import (
	"encoding/json"
	"strings"
	"time"

	"serica-go/internal/data"
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

type BookListQuery struct {
	Keyword     string `query:"keyword"`
	PageIndex   int    `query:"pageIndex,1"`
	PageSize    int    `query:"pageSize,20"`
	Category    string `query:"category"`
	FileType    string `query:"fileType"`
	PublishYear string `query:"publishYear"`
	Language    string `query:"language"`
	Order       string `query:"order"`
	Desc        bool   `query:"desc,true"`
}

func (q BookListQuery) Offset() int {
	return (q.PageIndex - 1) * q.PageSize
}

func parseCommaSep(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func (q BookListQuery) ToFilter() data.BookFilter {
	return data.BookFilter{
		Keyword:      q.Keyword,
		Category:     q.Category,
		FileTypes:    parseCommaSep(q.FileType),
		PublishYears: parseCommaSep(q.PublishYear),
		Languages:    parseCommaSep(q.Language),
		Order:        q.Order,
		Desc:         q.Desc,
	}
}

type PageQuery struct {
	PageIndex int `query:"page,1"`
	PageSize  int `query:"limit,20"`
}

func (q PageQuery) Offset() int {
	return (q.PageIndex - 1) * q.PageSize
}

type AuthorsQuery struct {
	Name      string `query:"name"`
	PageIndex int    `query:"page,1"`
	PageSize  int    `query:"limit,20"`
}

func (q AuthorsQuery) Offset() int {
	return (q.PageIndex - 1) * q.PageSize
}

type CategoriesPaginatedQuery struct {
	Keyword   string `query:"keyword"`
	PageIndex int    `query:"page,1"`
	PageSize  int    `query:"limit,20"`
}

func (q CategoriesPaginatedQuery) Offset() int {
	return (q.PageIndex - 1) * q.PageSize
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

type BookIDQuery struct {
	BookID int64 `query:"bookId,0"`
}

type ReadingPosQuery struct {
	BookID int64 `query:"bookId,0"`
}
