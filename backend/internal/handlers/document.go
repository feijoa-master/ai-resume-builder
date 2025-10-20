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

type DocumentHandler struct {
	documentService *service.DocumentService
}

func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

// GenerateResume generates a resume
func (h *DocumentHandler) GenerateResume(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req models.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Validate required fields
	if req.JobDescription == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_FIELDS", "Job description is required", nil)
		return
	}

	if req.TemplateID == "" {
		req.TemplateID = "classic" // Default template
	}

	// Set type
	req.Type = "resume"

	// Generate document
	doc, err := h.documentService.GenerateDocument(userID, &req)
	if err != nil {
		if err == service.ErrNoFreeGenerationsLeft {
			respondWithError(w, http.StatusForbidden, "NO_FREE_GENERATIONS", "No free generations left. Please upgrade to premium.", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "GENERATION_FAILED", "Failed to generate resume", nil)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"id":       doc.ID,
		"status":   "completed",
		"message":  "Resume generated successfully",
		"document": doc,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// GenerateCoverLetter generates a cover letter
func (h *DocumentHandler) GenerateCoverLetter(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req models.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Validate required fields
	if req.JobDescription == "" {
		respondWithError(w, http.StatusBadRequest, "MISSING_FIELDS", "Job description is required", nil)
		return
	}

	if req.TemplateID == "" {
		req.TemplateID = "classic"
	}

	// Set type
	req.Type = "cover_letter"

	// Generate document
	doc, err := h.documentService.GenerateDocument(userID, &req)
	if err != nil {
		if err == service.ErrNoFreeGenerationsLeft {
			respondWithError(w, http.StatusForbidden, "NO_FREE_GENERATIONS", "No free generations left. Please upgrade to premium.", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "GENERATION_FAILED", "Failed to generate cover letter", nil)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"id":       doc.ID,
		"status":   "completed",
		"message":  "Cover letter generated successfully",
		"document": doc,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// GetDocuments retrieves all user documents
func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	documents, err := h.documentService.GetDocuments(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get documents", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, documents)
}

// GetDocument retrieves a specific document
func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	vars := mux.Vars(r)
	docID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_ID", "Invalid document ID", nil)
		return
	}

	document, err := h.documentService.GetDocument(userID, docID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "NOT_FOUND", "Document not found", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, document)
}

// UpdateDocument updates a document
func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	vars := mux.Vars(r)
	docID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_ID", "Invalid document ID", nil)
		return
	}

	var doc models.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	doc.ID = docID
	doc.UserID = userID

	if err := h.documentService.UpdateDocument(userID, &doc); err != nil {
		respondWithError(w, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update document", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Document updated successfully"})
}

// DeleteDocument deletes a document
func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	vars := mux.Vars(r)
	docID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "INVALID_ID", "Invalid document ID", nil)
		return
	}

	if err := h.documentService.DeleteDocument(userID, docID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete document", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Document deleted successfully"})
}
