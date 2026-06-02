package exception

import "net/http"

const genericGroup = 100

var (
	InternalServerError = NewException(genericGroup, 1, "InternalServerException",
		map[string]string{
			"zh_HK": "服務內部錯誤",
			"zh_cn": "服务内部错误",
			"en_us": "Internal server error",
		}, http.StatusInternalServerError)

	InvalidBody = NewException(genericGroup, 2, "InvalidBodyException",
		map[string]string{
			"zh_HK": "請求參數無效",
			"zh_cn": "请求参数无效",
			"en_us": "Invalid request body",
		}, http.StatusBadRequest)

	Unauthorized = NewException(genericGroup, 3, "UnauthorizedException",
		map[string]string{
			"zh_HK": "未授權訪問",
			"zh_cn": "未授权访问",
			"en_us": "Unauthorized",
		}, http.StatusUnauthorized)
)
