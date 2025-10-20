package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrProfileNotFound   = errors.New("profile not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, free_generations_left, is_premium, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.FreeGenerationsLeft,
		user.IsPremium,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, free_generations_left, is_premium, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.FreeGenerationsLeft,
		&user.IsPremium,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, free_generations_left, is_premium, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.FreeGenerationsLeft,
		&user.IsPremium,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, free_generations_left = $3, is_premium = $4, updated_at = NOW()
		WHERE id = $5
	`

	result, err := r.db.Exec(query, user.Email, user.FullName, user.FreeGenerationsLeft, user.IsPremium, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// DecrementFreeGenerations decrements the free generations count
func (r *UserRepository) DecrementFreeGenerations(userID uuid.UUID) error {
	query := `
		UPDATE users
		SET free_generations_left = free_generations_left - 1, updated_at = NOW()
		WHERE id = $1 AND free_generations_left > 0
	`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to decrement free generations: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no free generations left or user not found")
	}

	return nil
}

// CreateProfile creates a profile for a user
func (r *UserRepository) CreateProfile(userID uuid.UUID) error {
	query := `
		INSERT INTO profiles (user_id, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
	`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}
