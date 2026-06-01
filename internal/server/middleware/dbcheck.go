package middleware

import (
	"encoding/json"
	"net/http"
)

func RequireDB(dbHealthy bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !dbHealthy {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":  503,
					"error": "database not available",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
