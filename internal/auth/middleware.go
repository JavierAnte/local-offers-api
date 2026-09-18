package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/JavierAnte/local-offers-api/internal/httpx"
	"github.com/google/uuid"
)

type contextKey string

const userIDContextKey contextKey = "userID"

// RequireAuth validates the Bearer token on the Authorization header and
// stores the authenticated user's id in the request context. Routes that
// don't wrap with this middleware stay public (browsing offers is public;
// only mutating endpoints require auth).
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				httpx.WriteError(w, http.StatusUnauthorized, "authentication_required", "A valid bearer token is required.")
				return
			}

			userID, err := ParseToken(jwtSecret, token)
			if err != nil {
				httpx.WriteError(w, http.StatusUnauthorized, "invalid_token", "The bearer token is invalid or expired.")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return userID, ok
}
