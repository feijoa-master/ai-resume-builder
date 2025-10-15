# AI Resume Builder - Backend

Go-based REST API for AI-powered resume and cover letter generation.

## 🛠 Tech Stack

- **Go 1.21+** - Programming language
- **PostgreSQL 15** - Database
- **Gorilla Mux** - HTTP router
- **JWT** - Authentication
- **OpenAI API** - AI generation
- **Docker** - Containerization

## 📁 Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration
│   ├── database/
│   │   └── database.go          # DB connection
│   ├── models/
│   │   └── models.go            # Data models
│   ├── utils/
│   │   ├── jwt.go               # JWT utilities
│   │   └── password.go          # Password utilities
│   ├── handlers/                # HTTP handlers (coming soon)
│   ├── middleware/              # HTTP middleware (coming soon)
│   ├── repository/              # Database repositories (coming soon)
│   └── service/                 # Business logic (coming soon)
├── migrations/
│   └── 001_initial_schema.sql   # Database schema
├── .env.example                 # Environment variables template
├── docker-compose.yml           # Docker configuration
├── Makefile                     # Build automation
├── go.mod                       # Go modules
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- PostgreSQL 15 (or use Docker)
- OpenAI API key

### 1. Clone and Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd backend

# Copy environment variables
cp .env.example .env

# Edit .env and add your OpenAI API key
nano .env
```

### 2. Start Database

```bash
# Start PostgreSQL with Docker
make docker-up

# Or manually
docker-compose up -d postgres
```

### 3. Run Migrations

```bash
# Apply database schema
make migrate-up

# Or manually
psql "postgres://postgres:postgres@localhost:5432/resume_builder" -f migrations/001_initial_schema.sql
```

### 4. Install Dependencies

```bash
make install

# Or manually
go mod download
```

### 5. Run the Application

```bash
# Development mode
make run

# Or manually
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## 🔧 Available Commands

```bash
make help          # Show all available commands
make install       # Install dependencies
make run           # Run the application
make build         # Build binary
make test          # Run tests
make docker-up     # Start Docker containers
make docker-down   # Stop Docker containers
make migrate-up    # Run migrations
make db-shell      # Open PostgreSQL shell
make lint          # Run linter
make fmt           # Format code
make dev           # Start full dev environment
```

## 📡 API Endpoints

### Health Check
```
GET /api/v1/health
```

### Authentication (Coming Soon)
```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

### Profile Management (Coming Soon)
```
GET    /api/v1/profile
PUT    /api/v1/profile
POST   /api/v1/profile/experience
PUT    /api/v1/profile/experience/:id
DELETE /api/v1/profile/experience/:id
```

### Document Generation (Coming Soon)
```
POST   /api/v1/generate/resume
POST   /api/v1/generate/cover-letter
GET    /api/v1/generate/status/:id
```

### Document Management (Coming Soon)
```
GET    /api/v1/documents
GET    /api/v1/documents/:id
PUT    /api/v1/documents/:id
DELETE /api/v1/documents/:id
GET    /api/v1/documents/:id/pdf
```

## 🔐 Environment Variables

See `.env.example` for all available configuration options.

Key variables:
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT tokens
- `OPENAI_API_KEY` - OpenAI API key for AI generation

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific test
go test -v ./internal/utils/...

# Run tests with coverage
go test -cover ./...
```

## 📦 Building

```bash
# Build binary
make build

# Run binary
./bin/resume-builder-api
```

## 🐳 Docker

```bash
# Build and run with Docker Compose
docker-compose up --build

# View logs
docker-compose logs -f backend

# Stop containers
docker-compose down
```

## 🔍 Database Access

```bash
# Open PostgreSQL shell
make db-shell

# Or manually
docker-compose exec postgres psql -U postgres -d resume_builder
```

## 📝 Development Workflow

1. Start database: `make docker-up`
2. Run migrations: `make migrate-up`
3. Start dev server: `make run`
4. Make changes and test
5. Format code: `make fmt`
6. Run tests: `make test`

## 🐛 Troubleshooting

### Database connection failed
```bash
# Check if PostgreSQL is running
docker-compose ps

# Check logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

### Port already in use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Migrations not applied
```bash
# Manually apply migrations
psql "postgres://postgres:postgres@localhost:5432/resume_builder" -f migrations/001_initial_schema.sql
```

## 📚 Next Steps

- [ ] Implement authentication handlers
- [ ] Add profile management
- [ ] Integrate OpenAI API
- [ ] Add document generation logic
- [ ] Implement PDF export
- [ ] Add tests
- [ ] Deploy to production

## 📄 License

MIT License - Feel free to use this project for learning and building!