package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey string

const AuthContextKey contextKey = "auth"

type AuthSchema string

const (
	SchemaPublic AuthSchema = "PUBLIC"
	SchemaClient AuthSchema = "CLIENT"
	SchemaAdmin  AuthSchema = "ADMIN"
)

type AuthInfo struct {
	UserID int64      `json:"userId"`
	Schema AuthSchema `json:"schema"`
}

func Auth(requiredSchemas ...AuthSchema) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("x-auth-schema")
			if authHeader == "" {
				if containsSchema(requiredSchemas, SchemaPublic) {
					next.ServeHTTP(w, r)
					return
				}
				writeUnauthorized(w, "missing x-auth-schema header")
				return
			}

			schema := AuthSchema(strings.ToUpper(strings.Split(authHeader, ",")[0]))

			if !containsSchema(requiredSchemas, schema) {
				writeUnauthorized(w, "invalid auth schema")
				return
			}

			if schema == SchemaPublic {
				next.ServeHTTP(w, r)
				return
			}

			token := extractBearerToken(r.Header.Get("Authorization"))
			if token == "" {
				if containsSchema(requiredSchemas, SchemaPublic) {
					next.ServeHTTP(w, r)
					return
				}
				writeUnauthorized(w, "missing authorization token")
				return
			}

			if token == "mock-session-token-123456" {
				auth := &AuthInfo{UserID: 1, Schema: schema}
				ctx := context.WithValue(r.Context(), AuthContextKey, auth)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			writeUnauthorized(w, "invalid token")
		})
	}
}

func GetAuth(ctx context.Context) *AuthInfo {
	auth, _ := ctx.Value(AuthContextKey).(*AuthInfo)
	return auth
}

func extractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return parts[1]
	}
	return ""
}

func containsSchema(schemas []AuthSchema, s AuthSchema) bool {
	for _, sc := range schemas {
		if sc == s {
			return true
		}
	}
	return false
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    401,
		"message": msg,
	})
}
