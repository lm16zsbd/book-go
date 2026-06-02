package bookskey

import (
	"net/http"

	"serica-go/internal/pkg/exception"
)

const booksKeyGroup = 201

var (
	KeyNotFound = exception.NewException(booksKeyGroup, 1, "KeyNotFoundException",
		map[string]string{
			"zh_HK": "密鑰不存在",
			"zh_cn": "密钥不存在",
			"en_us": "Key not found",
		}, http.StatusNotFound)
)
