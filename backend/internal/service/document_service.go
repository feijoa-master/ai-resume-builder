package service

import (
	"errors"
	"fmt"

	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/feijoa-master/ai-resume-builder/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrNoFreeGenerationsLeft = errors.New("no free generations left")
	ErrInvalidDocumentType   = errors.New("invalid document type")
)

type DocumentService struct {
	documentRepo  *repository.DocumentRepository
	profileRepo   *repository.ProfileRepository
	userRepo      *repository.UserRepository
	openaiService *OpenAIService
}

func NewDocumentService(
	documentRepo *repository.DocumentRepository,
	profileRepo *repository.ProfileRepository,
	userRepo *repository.UserRepository,
	openaiService *OpenAIService,
) *DocumentService {
	return &DocumentService{
		documentRepo:  documentRepo,
		profileRepo:   profileRepo,
		userRepo:      userRepo,
		openaiService: openaiService,
	}
}

// GenerateDocument generates a resume or cover letter
func (s *DocumentService) GenerateDocument(userID uuid.UUID, req *models.GenerateRequest) (*models.Document, error) {
	// Check if user has free generations left
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsPremium && user.FreeGenerationsLeft <= 0 {
		return nil, ErrNoFreeGenerationsLeft
	}

	// Get user profile data
	profileData, err := s.getProfileData(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile data: %w", err)
	}

	// Generate document using OpenAI
	var generated *GeneratedDocument
	switch req.Type {
	case "resume":
		generated, err = s.openaiService.GenerateResume(profileData, req.JobDescription)
	case "cover_letter":
		generated, err = s.openaiService.GenerateCoverLetter(profileData, req.JobDescription, req.CompanyName)
	default:
		return nil, ErrInvalidDocumentType
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate document: %w", err)
	}

	// Parse generated content
	var content map[string]interface{}
	content = map[string]interface{}{
		"raw_content": generated.Content,
	}

	// Create document record
	doc := &models.Document{
		ID:             uuid.New(),
		UserID:         userID,
		Type:           req.Type,
		Title:          s.generateTitle(req),
		Content:        content,
		TemplateID:     req.TemplateID,
		JobTitle:       req.JobTitle,
		CompanyName:    req.CompanyName,
		JobDescription: req.JobDescription,
		Status:         "final",
	}

	if err := s.documentRepo.CreateDocument(doc); err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	// Save generation history
	cost := s.calculateCost(generated.PromptTokens, generated.CompletionTokens)
	history := &models.GenerationHistory{
		ID:               uuid.New(),
		UserID:           userID,
		DocumentID:       doc.ID,
		PromptTokens:     generated.PromptTokens,
		CompletionTokens: generated.CompletionTokens,
		TotalCost:        cost,
		GenerationTimeMs: generated.GenerationTimeMs,
	}

	if err := s.documentRepo.CreateGenerationHistory(history); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to save generation history: %v\n", err)
	}

	// Decrement free generations if not premium
	if !user.IsPremium {
		if err := s.userRepo.DecrementFreeGenerations(userID); err != nil {
			// Log error but don't fail since document is already created
			fmt.Printf("Failed to decrement free generations: %v\n", err)
		}
	}

	return doc, nil
}

// GetDocument retrieves a document by ID
func (s *DocumentService) GetDocument(userID, docID uuid.UUID) (*models.Document, error) {
	return s.documentRepo.GetDocumentByID(docID, userID)
}

// GetDocuments retrieves all user documents
func (s *DocumentService) GetDocuments(userID uuid.UUID) ([]*models.Document, error) {
	return s.documentRepo.GetDocuments(userID)
}

// UpdateDocument updates a document
func (s *DocumentService) UpdateDocument(userID uuid.UUID, doc *models.Document) error {
	doc.UserID = userID
	return s.documentRepo.UpdateDocument(doc)
}

// DeleteDocument deletes a document
func (s *DocumentService) DeleteDocument(userID, docID uuid.UUID) error {
	return s.documentRepo.DeleteDocument(docID, userID)
}

// getProfileData gathers all profile data for generation
func (s *DocumentService) getProfileData(userID uuid.UUID) (*ProfileData, error) {
	// Get user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Get profile
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get experiences
	experiences, err := s.profileRepo.GetExperiences(profile.ID)
	if err != nil {
		return nil, err
	}

	// Get education
	education, err := s.profileRepo.GetEducation(profile.ID)
	if err != nil {
		return nil, err
	}

	// Get skills
	skills, err := s.profileRepo.GetSkills(profile.ID)
	if err != nil {
		return nil, err
	}

	return &ProfileData{
		FullName:    user.FullName,
		Email:       user.Email,
		Phone:       profile.Phone,
		Location:    profile.Location,
		LinkedInURL: profile.LinkedInURL,
		GithubURL:   profile.GithubURL,
		WebsiteURL:  profile.WebsiteURL,
		Summary:     profile.Summary,
		Experiences: experiences,
		Education:   education,
		Skills:      skills,
	}, nil
}

// generateTitle creates a title for the document
func (s *DocumentService) generateTitle(req *models.GenerateRequest) string {
	if req.CompanyName != "" && req.JobTitle != "" {
		return fmt.Sprintf("%s - %s", req.CompanyName, req.JobTitle)
	}
	if req.CompanyName != "" {
		return req.CompanyName
	}
	if req.JobTitle != "" {
		return req.JobTitle
	}
	if req.Type == "resume" {
		return "Resume"
	}
	return "Cover Letter"
}

// calculateCost estimates the cost of generation (simplified)
// Real pricing: https://openai.com/pricing
func (s *DocumentService) calculateCost(promptTokens, completionTokens int) float64 {
	// Example pricing for gpt-4o-mini (as of 2024):
	// $0.150 per 1M input tokens, $0.600 per 1M output tokens
	inputCost := float64(promptTokens) * 0.00000015
	outputCost := float64(completionTokens) * 0.0000006
	return inputCost + outputCost
}
