package middleware

import (
	"context"
	"net/http"
	"strings"

	//"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/feijoa-master/ai-resume-builder/internal/utils"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	IsPremiumKey contextKey = "is_premium"
)

// JWTAuth middleware validates JWT tokens
func JWTAuth(jwtManager *utils.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Authorization header must be Bearer token")
				return
			}

			token := parts[1]

			// Validate token
			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				if err == utils.ErrExpiredToken {
					respondWithError(w, http.StatusUnauthorized, "EXPIRED_TOKEN", "Token has expired")
					return
				}
				respondWithError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token")
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, IsPremiumKey, claims.IsPremium)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePremium middleware checks if user has premium subscription
func RequirePremium(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPremium, ok := r.Context().Value(IsPremiumKey).(bool)
		if !ok || !isPremium {
			respondWithError(w, http.StatusForbidden, "PREMIUM_REQUIRED", "This feature requires a premium subscription")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper functions

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetUserEmailFromContext extracts user email from request context
func GetUserEmailFromContext(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(UserEmailKey).(string)
	return email, ok
}

// IsPremiumUser checks if user is premium from request context
func IsPremiumUser(r *http.Request) bool {
	isPremium, ok := r.Context().Value(IsPremiumKey).(bool)
	return ok && isPremium
}

func respondWithError(w http.ResponseWriter, statusCode int, code, message string) {
	//errorResponse := models.ErrorResponse{
	//	Error: models.ErrorDetail{
	//		Code:    code,
	//		Message: message,
	//	},
	//}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error":{"code":"` + code + `","message":"` + message + `"}}`))
}
