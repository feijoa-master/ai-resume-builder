package service

import (
	"errors"
	"fmt"

	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/feijoa-master/ai-resume-builder/internal/repository"
	"github.com/feijoa-master/ai-resume-builder/internal/utils"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already registered")
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.TokenResponse, error) {
	// Validate password strength
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Check if email already exists
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil && err != repository.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:                  uuid.New(),
		Email:               req.Email,
		PasswordHash:        passwordHash,
		FullName:            req.FullName,
		FreeGenerationsLeft: 2, // Free tier: 2 generations
		IsPremium:           false,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		if err == repository.ErrUserAlreadyExists {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create empty profile for the user
	if err := s.userRepo.CreateProfile(user.ID); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.IsPremium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.IsPremium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Don't return password hash
	user.PasswordHash = ""

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *models.LoginRequest) (*models.TokenResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.IsPremium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.IsPremium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Don't return password hash
	user.PasswordHash = ""

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user from database to ensure they still exist
	user, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate new tokens
	newAccessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.IsPremium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, user.IsPremium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Don't return password hash
	user.PasswordHash = ""

	return &models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

// GetUserByID retrieves user information
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}
