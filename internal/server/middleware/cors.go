package middleware

import (
	"net/http"
	"strings"

	"serica-go/internal/conf"
)

func CORS(cfg conf.CORS) func(http.Handler) http.Handler {
	allowOrigin := "*"
	allowMethods := "GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS"
	allowHeaders := "Content-Type,Authorization,X-Auth-Schema,X-Device,X-Forwarded-For"

	if len(cfg.AllowedOrigins) > 0 {
		allowOrigin = strings.Join(cfg.AllowedOrigins, ",")
	}
	if len(cfg.AllowedMethods) > 0 {
		allowMethods = strings.Join(cfg.AllowedMethods, ",")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
