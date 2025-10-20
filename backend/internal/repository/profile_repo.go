package repository

import (
	"database/sql"
	"fmt"

	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// GetProfileByUserID retrieves a user's profile
func (r *ProfileRepository) GetProfileByUserID(userID uuid.UUID) (*models.Profile, error) {
	query := `
		SELECT id, user_id, phone, location, linkedin_url, github_url, website_url, summary, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	profile := &models.Profile{}

	// Use sql.NullString for nullable fields
	var phone, location, linkedinURL, githubURL, websiteURL, summary sql.NullString

	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&phone,
		&location,
		&linkedinURL,
		&githubURL,
		&websiteURL,
		&summary,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Convert sql.NullString to regular string
	profile.Phone = phone.String
	profile.Location = location.String
	profile.LinkedInURL = linkedinURL.String
	profile.GithubURL = githubURL.String
	profile.WebsiteURL = websiteURL.String
	profile.Summary = summary.String

	return profile, nil
}

// UpdateProfile updates profile information
func (r *ProfileRepository) UpdateProfile(profile *models.Profile) error {
	query := `
		UPDATE profiles
		SET phone = $1, location = $2, linkedin_url = $3, github_url = $4, 
		    website_url = $5, summary = $6, updated_at = NOW()
		WHERE user_id = $7
	`

	result, err := r.db.Exec(query,
		profile.Phone,
		profile.Location,
		profile.LinkedInURL,
		profile.GithubURL,
		profile.WebsiteURL,
		profile.Summary,
		profile.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
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

// Experience methods

func (r *ProfileRepository) CreateExperience(exp *models.Experience) error {
	query := `
		INSERT INTO experiences (id, profile_id, company, position, start_date, end_date, is_current, description, achievements, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		exp.ID,
		exp.ProfileID,
		exp.Company,
		exp.Position,
		exp.StartDate,
		exp.EndDate,
		exp.IsCurrent,
		exp.Description,
		pq.Array(exp.Achievements),
	).Scan(&exp.ID, &exp.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create experience: %w", err)
	}

	return nil
}

func (r *ProfileRepository) GetExperiences(profileID uuid.UUID) ([]*models.Experience, error) {
	query := `
		SELECT id, profile_id, company, position, start_date, end_date, is_current, description, achievements, created_at
		FROM experiences
		WHERE profile_id = $1
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experiences: %w", err)
	}
	defer rows.Close()

	var experiences []*models.Experience
	for rows.Next() {
		exp := &models.Experience{}
		err := rows.Scan(
			&exp.ID,
			&exp.ProfileID,
			&exp.Company,
			&exp.Position,
			&exp.StartDate,
			&exp.EndDate,
			&exp.IsCurrent,
			&exp.Description,
			pq.Array(&exp.Achievements),
			&exp.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan experience: %w", err)
		}
		experiences = append(experiences, exp)
	}

	return experiences, nil
}

func (r *ProfileRepository) UpdateExperience(exp *models.Experience) error {
	query := `
		UPDATE experiences
		SET company = $1, position = $2, start_date = $3, end_date = $4, 
		    is_current = $5, description = $6, achievements = $7
		WHERE id = $8 AND profile_id = $9
	`

	result, err := r.db.Exec(query,
		exp.Company,
		exp.Position,
		exp.StartDate,
		exp.EndDate,
		exp.IsCurrent,
		exp.Description,
		pq.Array(exp.Achievements),
		exp.ID,
		exp.ProfileID,
	)
	if err != nil {
		return fmt.Errorf("failed to update experience: %w", err)
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

func (r *ProfileRepository) DeleteExperience(id, profileID uuid.UUID) error {
	query := `DELETE FROM experiences WHERE id = $1 AND profile_id = $2`

	result, err := r.db.Exec(query, id, profileID)
	if err != nil {
		return fmt.Errorf("failed to delete experience: %w", err)
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

// Education methods

func (r *ProfileRepository) CreateEducation(edu *models.Education) error {
	query := `
		INSERT INTO education (id, profile_id, institution, degree, field_of_study, start_date, end_date, gpa, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		edu.ID,
		edu.ProfileID,
		edu.Institution,
		edu.Degree,
		edu.FieldOfStudy,
		edu.StartDate,
		edu.EndDate,
		edu.GPA,
	).Scan(&edu.ID, &edu.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create education: %w", err)
	}

	return nil
}

func (r *ProfileRepository) GetEducation(profileID uuid.UUID) ([]*models.Education, error) {
	query := `
		SELECT id, profile_id, institution, degree, field_of_study, start_date, end_date, gpa, created_at
		FROM education
		WHERE profile_id = $1
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education: %w", err)
	}
	defer rows.Close()

	var educations []*models.Education
	for rows.Next() {
		edu := &models.Education{}
		err := rows.Scan(
			&edu.ID,
			&edu.ProfileID,
			&edu.Institution,
			&edu.Degree,
			&edu.FieldOfStudy,
			&edu.StartDate,
			&edu.EndDate,
			&edu.GPA,
			&edu.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan education: %w", err)
		}
		educations = append(educations, edu)
	}

	return educations, nil
}

func (r *ProfileRepository) UpdateEducation(edu *models.Education) error {
	query := `
		UPDATE education
		SET institution = $1, degree = $2, field_of_study = $3, start_date = $4, 
		    end_date = $5, gpa = $6
		WHERE id = $7 AND profile_id = $8
	`

	result, err := r.db.Exec(query,
		edu.Institution,
		edu.Degree,
		edu.FieldOfStudy,
		edu.StartDate,
		edu.EndDate,
		edu.GPA,
		edu.ID,
		edu.ProfileID,
	)
	if err != nil {
		return fmt.Errorf("failed to update education: %w", err)
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

func (r *ProfileRepository) DeleteEducation(id, profileID uuid.UUID) error {
	query := `DELETE FROM education WHERE id = $1 AND profile_id = $2`

	result, err := r.db.Exec(query, id, profileID)
	if err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
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

// Skills methods

func (r *ProfileRepository) CreateSkill(skill *models.Skill) error {
	query := `
		INSERT INTO skills (id, profile_id, name, category, proficiency_level, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		skill.ID,
		skill.ProfileID,
		skill.Name,
		skill.Category,
		skill.ProficiencyLevel,
	).Scan(&skill.ID, &skill.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}

	return nil
}

func (r *ProfileRepository) GetSkills(profileID uuid.UUID) ([]*models.Skill, error) {
	query := `
		SELECT id, profile_id, name, category, proficiency_level, created_at
		FROM skills
		WHERE profile_id = $1
		ORDER BY category, name
	`

	rows, err := r.db.Query(query, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	defer rows.Close()

	var skills []*models.Skill
	for rows.Next() {
		skill := &models.Skill{}
		err := rows.Scan(
			&skill.ID,
			&skill.ProfileID,
			&skill.Name,
			&skill.Category,
			&skill.ProficiencyLevel,
			&skill.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (r *ProfileRepository) DeleteSkill(id, profileID uuid.UUID) error {
	query := `DELETE FROM skills WHERE id = $1 AND profile_id = $2`

	result, err := r.db.Exec(query, id, profileID)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
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
