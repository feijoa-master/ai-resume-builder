package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/feijoa-master/ai-resume-builder/internal/middleware"
	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/feijoa-master/ai-resume-builder/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile retrieves user's profile
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get profile", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, profile)
}

// UpdateProfile updates user's profile
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if err := h.profileService.UpdateProfile(userID, &profile); err != nil {
		respondWithError(w, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update profile", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

// Experience handlers

func (h *ProfileHandler) CreateExperience(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var exp models.Experience
	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Validate required fields
	if exp.Company == "" || exp.Position == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_FIELDS", "Company and position are required", nil)
		return
	}

	createdExp, err := h.profileService.CreateExperience(userID, &exp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create experience", nil)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdExp)
}

func (h *ProfileHandler) GetExperiences(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	experiences, err := h.profileService.GetExperiences(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get experiences", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, experiences)
}

func (h *ProfileHandler) UpdateExperience(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	vars := mux.Vars(r)
	expID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_ID", "Invalid experience ID", nil)
		return
	}

	var exp models.Experience
	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if err := h.profileService.UpdateExperience(userID, expID, &exp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update experience", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Experience updated successfully"})
}

func (h *ProfileHandler) DeleteExperience(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	vars := mux.Vars(r)
	expID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_ID", "Invalid experience ID", nil)
		return
	}

	if err := h.profileService.DeleteExperience(userID, expID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete experience", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Experience deleted successfully"})
}

// Education handlers

func (h *ProfileHandler) CreateEducation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var edu models.Education
	if err := json.NewDecoder(r.Body).Decode(&edu); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if edu.Institution == "" || edu.Degree == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_FIELDS", "Institution and degree are required", nil)
		return
	}

	createdEdu, err := h.profileService.CreateEducation(userID, &edu)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create education", nil)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdEdu)
}

func (h *ProfileHandler) GetEducation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	education, err := h.profileService.GetEducation(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get education", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, education)
}

func (h *ProfileHandler) UpdateEducation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	vars := mux.Vars(r)
	eduID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_ID", "Invalid education ID", nil)
		return
	}

	var edu models.Education
	if err := json.NewDecoder(r.Body).Decode(&edu); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if err := h.profileService.UpdateEducation(userID, eduID, &edu); err != nil {
		respondWithError(w, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update education", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Education updated successfully"})
}

func (h *ProfileHandler) DeleteEducation(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	vars := mux.Vars(r)
	eduID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_ID", "Invalid education ID", nil)
		return
	}

	if err := h.profileService.DeleteEducation(userID, eduID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete education", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Education deleted successfully"})
}

// Skills handlers

func (h *ProfileHandler) CreateSkill(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var skill models.Skill
	if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	if skill.Name == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_FIELDS", "Skill name is required", nil)
		return
	}

	createdSkill, err := h.profileService.CreateSkill(userID, &skill)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create skill", nil)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdSkill)
}

func (h *ProfileHandler) GetSkills(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	skills, err := h.profileService.GetSkills(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get skills", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, skills)
}

func (h *ProfileHandler) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	vars := mux.Vars(r)
	skillID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_ID", "Invalid skill ID", nil)
		return
	}

	if err := h.profileService.DeleteSkill(userID, skillID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete skill", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Skill deleted successfully"})
}
