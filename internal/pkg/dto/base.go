package dto

type PageResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	PageIndex  int   `json:"pageIndex"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

func NewPageResult[T any](items []T, total int64, pageIndex, pageSize int) PageResult[T] {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return PageResult[T]{
		Items:      items,
		Total:      total,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

func EmptyPageResult[T any](pageIndex, pageSize int) PageResult[T] {
	return PageResult[T]{
		Items:      []T{},
		Total:      0,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: 0,
	}
}
