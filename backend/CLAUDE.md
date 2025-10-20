# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI Resume Builder backend is a Go REST API that generates AI-powered resumes and cover letters using OpenAI. It features JWT authentication, PostgreSQL database, and a layered architecture (handlers → services → repositories).

## Development Commands

### Building and Running
```bash
make run              # Run development server (port 8080)
make build            # Build binary to bin/resume-builder-api
make dev              # Start database + run server

# Manual run
go run cmd/api/main.go
```

### Testing
```bash
make test             # Run all tests
go test -v ./...      # Run all tests with verbose output
go test -v ./internal/utils/...    # Run specific package tests
go test -cover ./...  # Run tests with coverage
```

### Database
```bash
make docker-up        # Start PostgreSQL container
make docker-down      # Stop containers
make migrate-up       # Apply schema (requires DATABASE_URL env)
make db-shell         # Open psql shell

# Manual migration
psql "postgres://postgres:postgres@localhost:5432/resume_builder" -f migrations/001_initial_schema.sql
```

### Code Quality
```bash
make fmt              # Format code with gofmt
make lint             # Run golangci-lint
make install          # Download and tidy dependencies
```

## Architecture

### Layered Structure
The codebase follows a clean architecture pattern with clear separation of concerns:

1. **cmd/api/main.go**: Application entry point. Initializes config, database, dependencies, and sets up HTTP routes with CORS and graceful shutdown.

2. **internal/handlers/**: HTTP request handlers. Parse requests, validate input, call service layer, format responses. Use standard error codes (INVALID_REQUEST, EMAIL_EXISTS, etc).

3. **internal/service/**: Business logic layer. Services orchestrate repositories and utilities. Example: `AuthService` handles registration (password validation, hash, create user, create empty profile, generate JWT tokens).

4. **internal/repository/**: Database access layer. Repositories execute SQL queries and map to models. Should return domain-specific errors like `ErrUserNotFound`, `ErrUserAlreadyExists`.

5. **internal/middleware/**: HTTP middleware (JWTAuth, RequirePremium). Auth middleware extracts Bearer token, validates JWT, and injects user context (UserIDKey, UserEmailKey, IsPremiumKey).

6. **internal/models/**: Data models and DTOs. Structs for database entities (User, Profile, Experience, Education, Skill, Document) and request/response DTOs (RegisterRequest, TokenResponse, etc).

7. **internal/utils/**: Shared utilities (JWT manager, password hashing).

8. **internal/config/**: Configuration loading from environment variables.

9. **internal/database/**: Database connection and health checks.

### Dependency Injection Flow
```
main.go creates: DB → Repositories → Services → Handlers
Router uses: Middleware (wraps handlers)
```

Dependencies are injected top-down, not created within layers.

### Authentication Flow
- Public routes: `/api/v1/auth/register`, `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/health`
- Protected routes: Use `middleware.JWTAuth(jwtManager)` and access user context via `middleware.GetUserIDFromContext(r)`
- JWT tokens: Access tokens expire in 15m, refresh tokens in 168h (7 days)
- Free tier: Users get 2 free generations, tracked in `users.free_generations_left`

### Database Schema
- **users**: Authentication and subscription info
- **profiles**: User professional data (unique per user, created on registration)
- **experiences/education/skills**: Profile sections (foreign key to profiles, cascade delete)
- **documents**: Generated resumes/cover letters (JSONB content field)
- **generation_history**: AI usage tracking and costs

Key relationships:
- User 1:1 Profile (created automatically on registration in `AuthService.Register`)
- Profile 1:N Experience, Education, Skill
- User 1:N Documents, GenerationHistory

### Error Handling Patterns
- Handlers return structured errors with code, message, and optional details
- Service layer returns Go errors (use `errors.New` or `fmt.Errorf` with `%w`)
- Repository layer returns domain errors (`ErrUserNotFound`, `ErrUserAlreadyExists`)
- Handlers map service/repo errors to HTTP status codes and error codes

### Context Usage
Middleware stores authenticated user data in request context:
```go
userID, ok := middleware.GetUserIDFromContext(r)
email, ok := middleware.GetUserEmailFromContext(r)
isPremium := middleware.IsPremiumUser(r)
```

## Key Technologies

- **gorilla/mux**: HTTP router with path variables (e.g., `{id}` in `/profile/experience/{id}`)
- **golang-jwt/jwt/v5**: JWT token generation and validation
- **lib/pq**: PostgreSQL driver
- **google/uuid**: UUID generation
- **rs/cors**: CORS handling for frontend origins
- **golang.org/x/crypto/bcrypt**: Password hashing

## Environment Configuration

Copy `.env.example` to `.env` and configure:
- `DATABASE_URL`: PostgreSQL connection string (default: localhost:5432/resume_builder)
- `JWT_SECRET`: Must be changed in production (min 32 chars recommended)
- `OPENAI_API_KEY`: Required for AI generation features
- `PORT`: Server port (default: 8080)
- `ALLOWED_ORIGINS`: Comma-separated CORS origins for frontend

## Testing Approach

When writing tests:
- Use table-driven tests for multiple scenarios
- Mock repositories/services with interfaces
- Test error cases and edge conditions
- Verify JWT token validation in middleware tests
- Test database constraints (unique emails, cascade deletes)

## Implementation Guidelines

1. **Adding new endpoints**: Follow the pattern in `cmd/api/main.go` (public vs protected routes)
2. **Database queries**: Use repositories with prepared statements, return domain errors
3. **Password handling**: Always use `utils.HashPassword` and `utils.CheckPassword`, never store plaintext
4. **Token generation**: Use `jwtManager.GenerateAccessToken` and `GenerateRefreshToken` with consistent claims
5. **Response formatting**: Use `respondWithJSON` and `respondWithError` helpers for consistency
6. **Validation**: Validate at handler layer before calling services

## Future Features (Not Yet Implemented)

The models and schema include structures for:
- Document generation endpoints (POST /generate/resume, /generate/cover-letter)
- PDF export functionality
- Generation history tracking
- Premium subscription features

When implementing these, follow the existing layered architecture pattern.
