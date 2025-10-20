package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/feijoa-master/ai-resume-builder/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email, password, and full name are required", nil)
		return
	}

	// Register user
	tokenResponse, err := h.authService.Register(&req)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			respondWithError(w, http.StatusConflict, "EMAIL_EXISTS", "Email already registered", nil)
			return
		}
		if strings.Contains(err.Error(), "password validation") {
			respondWithError(w, http.StatusBadRequest, "WEAK_PASSWORD", err.Error(), nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to register user", nil)
		return
	}

	respondWithJSON(w, http.StatusCreated, tokenResponse)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email and password are required", nil)
		return
	}

	// Authenticate user
	tokenResponse, err := h.authService.Login(&req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			respondWithError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "LOGIN_FAILED", "Failed to login", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, tokenResponse)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if req.RefreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_TOKEN", "Refresh token is required", nil)
		return
	}

	// Refresh tokens
	tokenResponse, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired refresh token", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, tokenResponse)
}

// Logout handles user logout (client-side token removal)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a JWT-based system, logout is typically handled client-side
	// by removing the tokens. We just acknowledge the request.
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// Helper functions for JSON responses

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, code, message string, details map[string]interface{}) {
	errorResponse := models.ErrorResponse{
		Error: models.ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	respondWithJSON(w, statusCode, errorResponse)
}
