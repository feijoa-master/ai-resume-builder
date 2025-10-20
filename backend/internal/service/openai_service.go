package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/feijoa-master/ai-resume-builder/internal/models"
	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
}

func NewOpenAIService(apiKey, model string, maxTokens int, temperature float64) *OpenAIService {
	client := openai.NewClient(apiKey)
	return &OpenAIService{
		client:      client,
		model:       model,
		maxTokens:   maxTokens,
		temperature: float32(temperature),
	}
}

// GenerateResume generates a resume based on profile and job description
func (s *OpenAIService) GenerateResume(profile *ProfileData, jobDescription string) (*GeneratedDocument, error) {
	prompt := s.buildResumePrompt(profile, jobDescription)

	startTime := time.Now()
	response, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       s.model,
			Temperature: s.temperature,
			MaxTokens:   s.maxTokens,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a professional resume writer with 10+ years of experience. Create ATS-friendly, impactful resumes that highlight candidates' strengths.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate resume: %w", err)
	}

	generationTime := time.Since(startTime).Milliseconds()

	// Parse the generated content
	content := response.Choices[0].Message.Content

	return &GeneratedDocument{
		Content:          content,
		PromptTokens:     response.Usage.PromptTokens,
		CompletionTokens: response.Usage.CompletionTokens,
		TotalTokens:      response.Usage.TotalTokens,
		GenerationTimeMs: int(generationTime),
		Model:            s.model,
	}, nil
}

// GenerateCoverLetter generates a cover letter based on profile and job description
func (s *OpenAIService) GenerateCoverLetter(profile *ProfileData, jobDescription, companyName string) (*GeneratedDocument, error) {
	prompt := s.buildCoverLetterPrompt(profile, jobDescription, companyName)

	startTime := time.Now()
	response, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       s.model,
			Temperature: s.temperature,
			MaxTokens:   800, // Cover letters are shorter
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert cover letter writer. Create compelling, personalized cover letters that showcase the candidate's fit for the role.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate cover letter: %w", err)
	}

	generationTime := time.Since(startTime).Milliseconds()

	content := response.Choices[0].Message.Content

	return &GeneratedDocument{
		Content:          content,
		PromptTokens:     response.Usage.PromptTokens,
		CompletionTokens: response.Usage.CompletionTokens,
		TotalTokens:      response.Usage.TotalTokens,
		GenerationTimeMs: int(generationTime),
		Model:            s.model,
	}, nil
}

// buildResumePrompt creates the prompt for resume generation
func (s *OpenAIService) buildResumePrompt(profile *ProfileData, jobDescription string) string {
	profileJSON, _ := json.MarshalIndent(profile, "", "  ")

	prompt := fmt.Sprintf(`Generate a professional, ATS-friendly resume based on the following candidate profile and job description.

CANDIDATE PROFILE:
%s

JOB DESCRIPTION:
%s

REQUIREMENTS:
1. Create a strong professional summary (3-4 sentences) that highlights key qualifications
2. List relevant work experience with bullet points focusing on achievements and impact
3. Include education and relevant skills
4. Use action verbs and quantify achievements where possible
5. Tailor the content to match the job requirements
6. Keep it concise and professional
7. Format in clean, readable sections

Return the resume in JSON format with the following structure:
{
  "summary": "Professional summary paragraph",
  "experience": [
    {
      "company": "Company Name",
      "position": "Job Title",
      "period": "Start - End",
      "highlights": ["Achievement 1", "Achievement 2", "Achievement 3"]
    }
  ],
  "education": [
    {
      "institution": "School Name",
      "degree": "Degree",
      "period": "Start - End"
    }
  ],
  "skills": {
    "technical": ["skill1", "skill2"],
    "soft": ["skill1", "skill2"]
  }
}`, profileJSON, jobDescription)

	return prompt
}

// buildCoverLetterPrompt creates the prompt for cover letter generation
func (s *OpenAIService) buildCoverLetterPrompt(profile *ProfileData, jobDescription, companyName string) string {
	profileJSON, _ := json.MarshalIndent(profile, "", "  ")

	prompt := fmt.Sprintf(`Generate a compelling cover letter based on the candidate profile and job description.

CANDIDATE PROFILE:
%s

JOB DESCRIPTION:
%s

COMPANY NAME: %s

REQUIREMENTS:
1. Start with a strong opening that shows enthusiasm and fit
2. Highlight 2-3 key qualifications that match the job requirements
3. Show understanding of the company/role
4. Demonstrate value the candidate brings
5. Close with a call to action
6. Keep it to 3-4 paragraphs
7. Professional but engaging tone

Return the cover letter in JSON format:
{
  "opening": "Opening paragraph",
  "body1": "First body paragraph",
  "body2": "Second body paragraph (if needed)",
  "closing": "Closing paragraph"
}`, profileJSON, jobDescription, companyName)

	return prompt
}

// ProfileData represents the complete user profile for generation
type ProfileData struct {
	FullName    string               `json:"full_name"`
	Email       string               `json:"email"`
	Phone       string               `json:"phone,omitempty"`
	Location    string               `json:"location,omitempty"`
	LinkedInURL string               `json:"linkedin_url,omitempty"`
	GithubURL   string               `json:"github_url,omitempty"`
	WebsiteURL  string               `json:"website_url,omitempty"`
	Summary     string               `json:"summary,omitempty"`
	Experiences []*models.Experience `json:"experiences"`
	Education   []*models.Education  `json:"education"`
	Skills      []*models.Skill      `json:"skills"`
}

// GeneratedDocument represents the result of AI generation
type GeneratedDocument struct {
	Content          string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	GenerationTimeMs int
	Model            string
}
