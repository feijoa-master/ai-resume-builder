package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/google/uuid"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// CreateDocument saves a generated document
func (r *DocumentRepository) CreateDocument(doc *models.Document) error {
	contentJSON, err := json.Marshal(doc.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	query := `
		INSERT INTO documents (id, user_id, type, title, content, template_id, job_title, company_name, job_description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRow(
		query,
		doc.ID,
		doc.UserID,
		doc.Type,
		doc.Title,
		contentJSON,
		doc.TemplateID,
		doc.JobTitle,
		doc.CompanyName,
		doc.JobDescription,
		doc.Status,
	).Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	return nil
}

// GetDocumentByID retrieves a document by ID
func (r *DocumentRepository) GetDocumentByID(id, userID uuid.UUID) (*models.Document, error) {
	query := `
		SELECT id, user_id, type, title, content, template_id, job_title, company_name, job_description, status, created_at, updated_at
		FROM documents
		WHERE id = $1 AND user_id = $2
	`

	doc := &models.Document{}
	var contentJSON []byte

	err := r.db.QueryRow(query, id, userID).Scan(
		&doc.ID,
		&doc.UserID,
		&doc.Type,
		&doc.Title,
		&contentJSON,
		&doc.TemplateID,
		&doc.JobTitle,
		&doc.CompanyName,
		&doc.JobDescription,
		&doc.Status,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Unmarshal content
	if err := json.Unmarshal(contentJSON, &doc.Content); err != nil {
		return nil, fmt.Errorf("failed to unmarshal content: %w", err)
	}

	return doc, nil
}

// GetDocuments retrieves all documents for a user
func (r *DocumentRepository) GetDocuments(userID uuid.UUID) ([]*models.Document, error) {
	query := `
		SELECT id, user_id, type, title, content, template_id, job_title, company_name, job_description, status, created_at, updated_at
		FROM documents
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		doc := &models.Document{}
		var contentJSON []byte

		err := rows.Scan(
			&doc.ID,
			&doc.UserID,
			&doc.Type,
			&doc.Title,
			&contentJSON,
			&doc.TemplateID,
			&doc.JobTitle,
			&doc.CompanyName,
			&doc.JobDescription,
			&doc.Status,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}

		// Unmarshal content
		if err := json.Unmarshal(contentJSON, &doc.Content); err != nil {
			return nil, fmt.Errorf("failed to unmarshal content: %w", err)
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// UpdateDocument updates a document
func (r *DocumentRepository) UpdateDocument(doc *models.Document) error {
	contentJSON, err := json.Marshal(doc.Content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	query := `
		UPDATE documents
		SET title = $1, content = $2, status = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
	`

	result, err := r.db.Exec(query, doc.Title, contentJSON, doc.Status, doc.ID, doc.UserID)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
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

// DeleteDocument deletes a document
func (r *DocumentRepository) DeleteDocument(id, userID uuid.UUID) error {
	query := `DELETE FROM documents WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
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

// CreateGenerationHistory saves generation metadata
func (r *DocumentRepository) CreateGenerationHistory(history *models.GenerationHistory) error {
	query := `
		INSERT INTO generation_history (id, user_id, document_id, prompt_tokens, completion_tokens, total_cost, generation_time_ms, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`

	_, err := r.db.Exec(
		query,
		history.ID,
		history.UserID,
		history.DocumentID,
		history.PromptTokens,
		history.CompletionTokens,
		history.TotalCost,
		history.GenerationTimeMs,
	)

	if err != nil {
		return fmt.Errorf("failed to create generation history: %w", err)
	}

	return nil
}
