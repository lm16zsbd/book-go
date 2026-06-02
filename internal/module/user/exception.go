package user

import (
	"net/http"

	"serica-go/internal/pkg/exception"
)

const userGroup = 220

var (
	UserNotFound = exception.NewException(userGroup, 1, "UserNotFoundException",
		map[string]string{
			"zh_HK": "用戶不存在",
			"zh_cn": "用户不存在",
			"en_us": "User not found",
		}, http.StatusNotFound)
)
