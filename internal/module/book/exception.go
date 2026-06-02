package book

import (
	"net/http"

	"serica-go/internal/pkg/exception"
)

const bookGroup = 200

var (
	BookNotFound = exception.NewException(bookGroup, 1, "BookNotFoundException",
		map[string]string{
			"zh_HK": "書籍不存在",
			"zh_cn": "书籍不存在",
			"en_us": "Book not found",
		}, http.StatusNotFound)

	CategoryNotFound = exception.NewException(bookGroup, 2, "CategoryNotFoundException",
		map[string]string{
			"zh_HK": "分類不存在",
			"zh_cn": "分类不存在",
			"en_us": "Category not found",
		}, http.StatusNotFound)
)
