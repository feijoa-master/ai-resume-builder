package service

import (
	"fmt"
	"log"

	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/feijoa-master/ai-resume-builder/internal/repository"
	"github.com/google/uuid"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

// GetProfile retrieves user's profile with all related data
func (s *ProfileService) GetProfile(userID uuid.UUID) (*models.Profile, error) {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		log.Printf("❌ GetProfile SQL error for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return profile, nil
}

// UpdateProfile updates profile information
func (s *ProfileService) UpdateProfile(userID uuid.UUID, profile *models.Profile) error {
	profile.UserID = userID

	if err := s.profileRepo.UpdateProfile(profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// Experience methods

func (s *ProfileService) CreateExperience(userID uuid.UUID, exp *models.Experience) (*models.Experience, error) {
	// Get profile to verify ownership
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	exp.ID = uuid.New()
	exp.ProfileID = profile.ID

	if err := s.profileRepo.CreateExperience(exp); err != nil {
		return nil, fmt.Errorf("failed to create experience: %w", err)
	}

	return exp, nil
}

func (s *ProfileService) GetExperiences(userID uuid.UUID) ([]*models.Experience, error) {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	experiences, err := s.profileRepo.GetExperiences(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experiences: %w", err)
	}

	return experiences, nil
}

func (s *ProfileService) UpdateExperience(userID uuid.UUID, expID uuid.UUID, exp *models.Experience) error {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	exp.ID = expID
	exp.ProfileID = profile.ID

	if err := s.profileRepo.UpdateExperience(exp); err != nil {
		return fmt.Errorf("failed to update experience: %w", err)
	}

	return nil
}

func (s *ProfileService) DeleteExperience(userID uuid.UUID, expID uuid.UUID) error {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	if err := s.profileRepo.DeleteExperience(expID, profile.ID); err != nil {
		return fmt.Errorf("failed to delete experience: %w", err)
	}

	return nil
}

// Education methods

func (s *ProfileService) CreateEducation(userID uuid.UUID, edu *models.Education) (*models.Education, error) {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	edu.ID = uuid.New()
	edu.ProfileID = profile.ID

	if err := s.profileRepo.CreateEducation(edu); err != nil {
		return nil, fmt.Errorf("failed to create education: %w", err)
	}

	return edu, nil
}

func (s *ProfileService) GetEducation(userID uuid.UUID) ([]*models.Education, error) {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	education, err := s.profileRepo.GetEducation(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education: %w", err)
	}

	return education, nil
}

func (s *ProfileService) UpdateEducation(userID uuid.UUID, eduID uuid.UUID, edu *models.Education) error {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	edu.ID = eduID
	edu.ProfileID = profile.ID

	if err := s.profileRepo.UpdateEducation(edu); err != nil {
		return fmt.Errorf("failed to update education: %w", err)
	}

	return nil
}

func (s *ProfileService) DeleteEducation(userID uuid.UUID, eduID uuid.UUID) error {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	if err := s.profileRepo.DeleteEducation(eduID, profile.ID); err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}

	return nil
}

// Skills methods

func (s *ProfileService) CreateSkill(userID uuid.UUID, skill *models.Skill) (*models.Skill, error) {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	skill.ID = uuid.New()
	skill.ProfileID = profile.ID

	if err := s.profileRepo.CreateSkill(skill); err != nil {
		return nil, fmt.Errorf("failed to create skill: %w", err)
	}

	return skill, nil
}

func (s *ProfileService) GetSkills(userID uuid.UUID) ([]*models.Skill, error) {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	skills, err := s.profileRepo.GetSkills(profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	return skills, nil
}

func (s *ProfileService) DeleteSkill(userID uuid.UUID, skillID uuid.UUID) error {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	if err := s.profileRepo.DeleteSkill(skillID, profile.ID); err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	return nil
}
