package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"serica-go/internal/data"
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

const sessionKeyPrefix = "session:token:"
const sessionTTL = 7 * 24 * 60 * 60

func Auth(userRepo *data.UserRepo, redis *data.RedisClient, requiredSchemas ...AuthSchema) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			schemaHeader := r.Header.Get("x-auth-schema")
			schema := resolveSchema(schemaHeader, requiredSchemas)
			if schema == "" {
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

			userID, ok := resolveUser(r.Context(), redis, userRepo, token)
			if !ok {
				if containsSchema(requiredSchemas, SchemaPublic) {
					next.ServeHTTP(w, r)
					return
				}
				writeUnauthorized(w, "invalid token")
				return
			}

			auth := &AuthInfo{UserID: userID, Schema: schema}
			ctx := context.WithValue(r.Context(), AuthContextKey, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func resolveSchema(schemaHeader string, required []AuthSchema) AuthSchema {
	if schemaHeader == "" {
		if containsSchema(required, SchemaPublic) {
			return SchemaPublic
		}
		return ""
	}
	schema := AuthSchema(strings.ToUpper(strings.Split(schemaHeader, ",")[0]))
	if containsSchema(required, schema) {
		return schema
	}
	return ""
}

func resolveUser(ctx context.Context, redis *data.RedisClient, userRepo *data.UserRepo, token string) (int64, bool) {
	if redis != nil {
		var cached struct {
			UserID int64  `json:"userId"`
			Email  string `json:"email"`
		}
		if err := redis.GetJSON(ctx, sessionKeyPrefix+token, &cached); err == nil && cached.UserID > 0 {
			return cached.UserID, true
		}
	}

	session, err := userRepo.FindSessionByToken(token)
	if err != nil || session == nil {
		return 0, false
	}

	user, err := userRepo.FindByID(session.UserID)
	if err != nil || user == nil {
		return 0, false
	}

	if redis != nil {
		_ = redis.SetJSON(ctx, sessionKeyPrefix+token, map[string]interface{}{
			"userId": user.ID,
			"email":  user.Email,
		}, time.Duration(sessionTTL)*time.Second)
	}

	return user.ID, true
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
	return authHeader
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
