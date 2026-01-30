package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/clearstatus/backend/internal/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserContext struct {
	ID    string
	Email string
}

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"invalid authorization header"}`, http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]
			claims, err := auth.ValidateJWT(secret, tokenString)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, &UserContext{ID: claims.UserID, Email: claims.Email})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) *UserContext {
	u, _ := ctx.Value(UserContextKey).(*UserContext)
	return u
}
