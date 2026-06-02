package middleware

import (
	"encoding/json"
	"net/http"
)

func NotFound(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// After the chain runs, if status is still 0 (no handler matched),
		// the actual 404 was sent but we let chi's built-in handle it.
		// This middleware is a passthrough for logging/metrics if needed.
	})
}

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    404,
			"message": "route not found",
			"data": map[string]interface{}{
				"name":                "RouteNotFoundException",
				"multilingualMessage": map[string]string{"zh_HK": "路由不存在", "zh_cn": "路由不存在", "en_us": "Route not found"},
				"params":              map[string]interface{}{},
			},
		})
	}
}
