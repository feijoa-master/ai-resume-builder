package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user
type User struct {
	ID                  uuid.UUID `json:"id"`
	Email               string    `json:"email"`
	PasswordHash        string    `json:"-"` // Never expose in JSON
	FullName            string    `json:"full_name"`
	FreeGenerationsLeft int       `json:"free_generations_left"`
	IsPremium           bool      `json:"is_premium"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Profile represents user's professional profile
type Profile struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Phone       string    `json:"phone,omitempty"`
	Location    string    `json:"location,omitempty"`
	LinkedInURL string    `json:"linkedin_url,omitempty"`
	GithubURL   string    `json:"github_url,omitempty"`
	WebsiteURL  string    `json:"website_url,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Experience represents work experience
type Experience struct {
	ID           uuid.UUID `json:"id"`
	ProfileID    uuid.UUID `json:"profile_id"`
	Company      string    `json:"company"`
	Position     string    `json:"position"`
	StartDate    Date      `json:"start_date"`
	EndDate      NullDate  `json:"end_date,omitempty"`
	IsCurrent    bool      `json:"is_current"`
	Description  string    `json:"description,omitempty"`
	Achievements []string  `json:"achievements,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// Education represents educational background
type Education struct {
	ID           uuid.UUID `json:"id"`
	ProfileID    uuid.UUID `json:"profile_id"`
	Institution  string    `json:"institution"`
	Degree       string    `json:"degree"`
	FieldOfStudy string    `json:"field_of_study,omitempty"`
	StartDate    Date      `json:"start_date"`
	EndDate      NullDate  `json:"end_date,omitempty"`
	GPA          float64   `json:"gpa,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// Skill represents a user skill
type Skill struct {
	ID               uuid.UUID `json:"id"`
	ProfileID        uuid.UUID `json:"profile_id"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`                    // technical, soft, language
	ProficiencyLevel string    `json:"proficiency_level,omitempty"` // beginner, intermediate, advanced, expert
	CreatedAt        time.Time `json:"created_at"`
}

// Document represents generated resume or cover letter
type Document struct {
	ID             uuid.UUID              `json:"id"`
	UserID         uuid.UUID              `json:"user_id"`
	Type           string                 `json:"type"` // resume, cover_letter
	Title          string                 `json:"title"`
	Content        map[string]interface{} `json:"content"` // JSON content
	TemplateID     string                 `json:"template_id,omitempty"`
	JobTitle       string                 `json:"job_title,omitempty"`
	CompanyName    string                 `json:"company_name,omitempty"`
	JobDescription string                 `json:"job_description,omitempty"`
	Status         string                 `json:"status"` // draft, final
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// GenerationHistory tracks AI generation usage
type GenerationHistory struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	DocumentID       uuid.UUID `json:"document_id"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalCost        float64   `json:"total_cost"`
	GenerationTimeMs int       `json:"generation_time_ms"`
	CreatedAt        time.Time `json:"created_at"`
}

// DTOs (Data Transfer Objects)

// RegisterRequest for user registration
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=2"`
}

// LoginRequest for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse contains JWT tokens
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

// GenerateRequest for document generation
type GenerateRequest struct {
	Type           string   `json:"type" validate:"required,oneof=resume cover_letter"`
	JobDescription string   `json:"job_description" validate:"required"`
	JobTitle       string   `json:"job_title,omitempty"`
	CompanyName    string   `json:"company_name,omitempty"`
	TemplateID     string   `json:"template_id" validate:"required"`
	CustomSections []string `json:"custom_sections,omitempty"`
}

// GenerateResponse after starting generation
type GenerateResponse struct {
	ID            uuid.UUID `json:"id"`
	Status        string    `json:"status"`
	EstimatedTime int       `json:"estimated_time"` // in seconds
}

// ErrorResponse for API errors
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}
