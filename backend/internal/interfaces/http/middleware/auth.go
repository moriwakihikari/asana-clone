package middleware

import (
	"context"
	"net/http"
	"strings"

	"asana-clone-backend/internal/infrastructure/auth"
	httpErrors "asana-clone-backend/internal/interfaces/errors"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware returns an HTTP middleware that validates Bearer tokens.
func AuthMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpErrors.RespondWithJSON(w, http.StatusUnauthorized, httpErrors.ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "missing authorization header",
				})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				httpErrors.RespondWithJSON(w, http.StatusUnauthorized, httpErrors.ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "invalid authorization header format",
				})
				return
			}

			tokenStr := parts[1]
			claims, err := jwtService.ValidateAccessToken(tokenStr)
			if err != nil {
				httpErrors.RespondWithJSON(w, http.StatusUnauthorized, httpErrors.ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "invalid or expired token",
				})
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the authenticated user's ID from the request context.
func GetUserID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(userIDKey).(uuid.UUID)
	return id
}
