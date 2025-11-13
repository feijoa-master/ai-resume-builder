# 📚 AI Resume Builder - Полная документация (Обучающий режим)

## 📖 Содержание

1. [Введение](#введение)
2. [Общая архитектура системы](#общая-архитектура-системы)
3. [Backend (Go)](#backend-go)
4. [Frontend (React + TypeScript)](#frontend-react--typescript)
5. [База данных](#база-данных)
6. [Аутентификация и безопасность](#аутентификация-и-безопасность)
7. [Интеграция с OpenAI](#интеграция-с-openai)
8. [Разработка и тестирование](#разработка-и-тестирование)
9. [Deployment](#deployment)
10. [Типичные задачи](#типичные-задачи)

---

## Введение

### Что это за проект?

**AI Resume Builder** - это веб-приложение для автоматической генерации резюме и сопроводительных писем с помощью искусственного интеллекта (OpenAI API).

### Основные возможности:

- Регистрация и аутентификация пользователей (JWT)
- Управление профилем (опыт работы, образование, навыки)
- AI-генерация резюме и cover letters
- Шаблоны документов
- История генерации
- Экспорт в PDF
- Бесплатный tier (2 генерации) + Premium подписка

### Стек технологий:

**Backend:**
- Go 1.21+ (язык программирования)
- Gorilla Mux (HTTP роутер)
- PostgreSQL 15 (база данных)
- JWT (аутентификация)
- OpenAI API (AI генерация)
- Docker (контейнеризация)

**Frontend:**
- React 18
- TypeScript
- Vite (сборка)
- TailwindCSS (стили)
- Zustand (state management)
- React Router (навигация)

---

## Общая архитектура системы

### Высокоуровневая схема

```
┌─────────────────────────────────────────────────────────────┐
│                        ПОЛЬЗОВАТЕЛЬ                          │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  FRONTEND (React + TS)                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Auth Pages  │  │   Dashboard  │  │   Profile    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Generate   │  │   History    │  │   Settings   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                              │
│  📦 State: Zustand Store (auth, profile, documents)         │
│  🔌 API Client: Axios + Interceptors                         │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP/REST
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   BACKEND (Go REST API)                      │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Handlers (HTTP Layer)                              │    │
│  │  ├── auth.go         (Register, Login, Refresh)    │    │
│  │  ├── profile.go      (CRUD Profile, Experience)    │    │
│  │  ├── document.go     (Generate, List, Export)      │    │
│  └──────────────────────────┬──────────────────────────────┘│
│                             │                               │
│  ┌──────────────────────────▼──────────────────────────┐   │
│  │  Middleware                                          │   │
│  │  ├── JWT Authentication                             │   │
│  │  ├── CORS                                           │   │
│  │  ├── Require Premium                                │   │
│  └──────────────────────────┬──────────────────────────────┘│
│                             │                               │
│  ┌──────────────────────────▼──────────────────────────┐   │
│  │  Services (Business Logic)                          │   │
│  │  ├── AuthService      (Password, JWT, Register)    │   │
│  │  ├── ProfileService   (Validation, Orchestration)  │   │
│  │  ├── DocumentService  (Generation logic)           │   │
│  │  ├── OpenAIService    (AI Integration)             │   │
│  └──────────────────────────┬──────────────────────────────┘│
│                             │                               │
│  ┌──────────────────────────▼──────────────────────────┐   │
│  │  Repositories (Data Access)                         │   │
│  │  ├── UserRepo         (SQL queries for users)      │   │
│  │  ├── ProfileRepo      (SQL queries for profiles)   │   │
│  │  ├── DocumentRepo     (SQL queries for docs)       │   │
│  └──────────────────────────┬──────────────────────────────┘│
│                             │                               │
└─────────────────────────────┼───────────────────────────────┘
                              │ SQL
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    PostgreSQL Database                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    users     │  │   profiles   │  │  experiences │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  education   │  │    skills    │  │  documents   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘

                              │ HTTP
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      OpenAI API                              │
│  - GPT-4 / GPT-3.5-turbo                                     │
│  - Генерация контента для резюме                             │
└─────────────────────────────────────────────────────────────┘
```

### Принципы архитектуры

1. **Разделение ответственности (Separation of Concerns)**
   - Backend и Frontend - отдельные приложения
   - Чистая архитектура с слоями (Clean Architecture)

2. **Dependency Injection (Внедрение зависимостей)**
   - Зависимости создаются в main.go и передаются вниз
   - Не создаем зависимости внутри слоев

3. **RESTful API**
   - Стандартные HTTP методы (GET, POST, PUT, DELETE)
   - Структурированные JSON responses
   - Коды ошибок и статусы

---

## Backend (Go)

### Структура проекта

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # 🚪 Entry point приложения
├── internal/
│   ├── config/
│   │   └── config.go            # ⚙️  Конфигурация из .env
│   ├── database/
│   │   └── database.go          # 🗄️  Подключение к БД
│   ├── models/
│   │   ├── models.go            # 📦 Структуры данных
│   │   └── date.go              # 📅 Кастомный Date type
│   ├── handlers/                # 🌐 HTTP обработчики
│   │   ├── auth.go
│   │   ├── profile.go
│   │   └── document.go
│   ├── middleware/              # 🔒 Middleware
│   │   └── auth.go
│   ├── service/                 # 🧠 Бизнес-логика
│   │   ├── auth_service.go
│   │   ├── profile_service.go
│   │   ├── document_service.go
│   │   └── openai_service.go
│   ├── repository/              # 📊 Доступ к БД
│   │   ├── user_repo.go
│   │   ├── profile_repo.go
│   │   └── document_repo.go
│   └── utils/                   # 🛠️  Утилиты
│       ├── jwt.go
│       └── password.go
├── migrations/
│   └── 001_initial_schema.sql   # 🗃️  Схема БД
├── qa/                          # 🧪 QA тесты
├── .env                         # 🔐 Переменные окружения
├── .env.example                 # 📝 Шаблон .env
├── docker-compose.yml           # 🐳 Docker конфигурация
├── Makefile                     # ⚡ Команды для разработки
├── go.mod                       # 📦 Go модули
└── README.md                    # 📖 Документация
```

### Подробный разбор слоев

#### 1. Entry Point (cmd/api/main.go)

**Что происходит при запуске:**

```go
func main() {
    // 1️⃣ Загружаем конфигурацию из .env
    cfg := config.LoadConfig()

    // 2️⃣ Подключаемся к базе данных
    db := database.Connect(cfg.DatabaseURL)
    defer db.Close()

    // 3️⃣ Создаем зависимости (Dependency Injection)
    jwtManager := utils.NewJWTManager(cfg.JWTSecret)

    // Repositories (доступ к БД)
    userRepo := repository.NewUserRepository(db)
    profileRepo := repository.NewProfileRepository(db)
    documentRepo := repository.NewDocumentRepository(db)

    // Services (бизнес-логика)
    authService := service.NewAuthService(userRepo, profileRepo, jwtManager)
    profileService := service.NewProfileService(profileRepo)
    openAIService := service.NewOpenAIService(cfg.OpenAIKey)
    documentService := service.NewDocumentService(
        documentRepo,
        profileRepo,
        openAIService,
    )

    // Handlers (HTTP обработчики)
    authHandler := handlers.NewAuthHandler(authService, jwtManager)
    profileHandler := handlers.NewProfileHandler(profileService)
    documentHandler := handlers.NewDocumentHandler(documentService)

    // 4️⃣ Настраиваем роутер (маршрутизация)
    router := mux.NewRouter()

    // Публичные роуты (не требуют авторизации)
    router.HandleFunc("/api/v1/health", healthHandler).Methods("GET")
    router.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods("POST")
    router.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")
    router.HandleFunc("/api/v1/auth/refresh", authHandler.Refresh).Methods("POST")

    // Защищенные роуты (требуют JWT токен)
    protected := router.PathPrefix("/api/v1").Subrouter()
    protected.Use(middleware.JWTAuth(jwtManager)) // Middleware для проверки токена

    protected.HandleFunc("/profile", profileHandler.GetProfile).Methods("GET")
    protected.HandleFunc("/profile", profileHandler.UpdateProfile).Methods("PUT")
    protected.HandleFunc("/profile/experience", profileHandler.CreateExperience).Methods("POST")
    // ... другие роуты

    // 5️⃣ Настраиваем CORS (для frontend)
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:5173"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
        AllowCredentials: true,
    }).Handler(router)

    // 6️⃣ Запускаем HTTP сервер
    server := &http.Server{
        Addr:    ":8080",
        Handler: corsHandler,
    }

    log.Println("Server started on :8080")
    server.ListenAndServe()
}
```

**Ключевые концепции:**

- **Dependency Injection**: Создаем все зависимости в одном месте и передаем их в конструкторы
- **Инициализация сверху вниз**: DB → Repositories → Services → Handlers
- **Разделение роутов**: Публичные и защищенные (с JWT middleware)

#### 2. Handlers (HTTP Layer)

**Задача:** Обработка HTTP запросов и ответов

**Паттерн работы handler'а:**

```go
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1️⃣ Парсим JSON из request body
    var req models.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
        return
    }

    // 2️⃣ Валидация входных данных
    if !isValidEmail(req.Email) {
        respondWithError(w, http.StatusBadRequest, "INVALID_EMAIL", "Invalid email format")
        return
    }
    if len(req.Password) < 8 {
        respondWithError(w, http.StatusBadRequest, "WEAK_PASSWORD", "Password must be at least 8 characters")
        return
    }

    // 3️⃣ Вызываем service layer (бизнес-логика)
    user, tokens, err := h.authService.Register(r.Context(), req)
    if err != nil {
        // Обработка ошибок из service layer
        if err == service.ErrEmailExists {
            respondWithError(w, http.StatusConflict, "EMAIL_EXISTS", "User already exists")
            return
        }
        respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register")
        return
    }

    // 4️⃣ Формируем успешный ответ
    response := models.TokenResponse{
        AccessToken:  tokens.AccessToken,
        RefreshToken: tokens.RefreshToken,
        User:         user,
    }

    // 5️⃣ Отправляем JSON ответ
    respondWithJSON(w, http.StatusCreated, response)
}
```

**Принципы:**

- Handler **не содержит** бизнес-логику
- Handler **валидирует** входные данные
- Handler **вызывает** service layer
- Handler **форматирует** ответы (JSON)
- Handler **обрабатывает** HTTP статусы и коды ошибок

**Стандартные коды ошибок:**

```go
const (
    ErrorCodeInvalidRequest  = "INVALID_REQUEST"
    ErrorCodeUnauthorized    = "UNAUTHORIZED"
    ErrorCodeForbidden       = "FORBIDDEN"
    ErrorCodeNotFound        = "NOT_FOUND"
    ErrorCodeEmailExists     = "EMAIL_EXISTS"
    ErrorCodeInternalError   = "INTERNAL_ERROR"
)
```

#### 3. Services (Business Logic Layer)

**Задача:** Реализация бизнес-логики приложения

**Пример: AuthService.Register**

```go
func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, *Tokens, error) {
    // 1️⃣ Проверяем, существует ли пользователь
    existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, nil, ErrEmailExists
    }

    // 2️⃣ Хешируем пароль (НИКОГДА не храним пароли в открытом виде!)
    passwordHash, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to hash password: %w", err)
    }

    // 3️⃣ Создаем пользователя
    user := &models.User{
        ID:                  uuid.New(),
        Email:               req.Email,
        PasswordHash:        passwordHash,
        FullName:            req.FullName,
        FreeGenerationsLeft: 2, // Бесплатные генерации
        IsPremium:           false,
        CreatedAt:           time.Now(),
        UpdatedAt:           time.Now(),
    }

    // 4️⃣ Сохраняем пользователя в БД
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, nil, fmt.Errorf("failed to create user: %w", err)
    }

    // 5️⃣ Создаем пустой профиль (связь 1:1 с пользователем)
    profile := &models.Profile{
        ID:        uuid.New(),
        UserID:    user.ID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := s.profileRepo.Create(ctx, profile); err != nil {
        // Если создание профиля провалилось, откатываем создание пользователя
        // В идеале это должно быть в транзакции
        return nil, nil, fmt.Errorf("failed to create profile: %w", err)
    }

    // 6️⃣ Генерируем JWT токены
    accessToken, err := s.jwtManager.GenerateAccessToken(user)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
    }

    refreshToken, err := s.jwtManager.GenerateRefreshToken(user)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
    }

    tokens := &Tokens{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }

    // 7️⃣ Возвращаем результат
    return user, tokens, nil
}
```

**Принципы:**

- Service **координирует** работу нескольких repositories
- Service **реализует** бизнес-правила (например, 2 бесплатных генерации)
- Service **обрабатывает** транзакции (в идеале)
- Service **возвращает** domain errors (не HTTP коды!)

**Domain Errors:**

```go
var (
    ErrEmailExists      = errors.New("email already exists")
    ErrInvalidPassword  = errors.New("invalid password")
    ErrUserNotFound     = errors.New("user not found")
    ErrInsufficientCredits = errors.New("insufficient generation credits")
)
```

#### 4. Repositories (Data Access Layer)

**Задача:** Выполнение SQL запросов к базе данных

**Пример: UserRepository**

```go
type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Создать пользователя
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (id, email, password_hash, full_name, free_generations_left, is_premium, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

    _, err := r.db.ExecContext(ctx, query,
        user.ID,
        user.Email,
        user.PasswordHash,
        user.FullName,
        user.FreeGenerationsLeft,
        user.IsPremium,
        user.CreatedAt,
        user.UpdatedAt,
    )

    if err != nil {
        // Проверяем constraint ошибки (например, unique email)
        if strings.Contains(err.Error(), "duplicate key") {
            return ErrUserAlreadyExists
        }
        return fmt.Errorf("failed to insert user: %w", err)
    }

    return nil
}

// Получить пользователя по email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `
        SELECT id, email, password_hash, full_name, free_generations_left, is_premium, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    user := &models.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.FullName,
        &user.FreeGenerationsLeft,
        &user.IsPremium,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    return user, nil
}

// Обновить количество бесплатных генераций
func (r *UserRepository) DecrementFreeGenerations(ctx context.Context, userID uuid.UUID) error {
    query := `
        UPDATE users
        SET free_generations_left = free_generations_left - 1,
            updated_at = $2
        WHERE id = $1 AND free_generations_left > 0
    `

    result, err := r.db.ExecContext(ctx, query, userID, time.Now())
    if err != nil {
        return fmt.Errorf("failed to decrement generations: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return ErrInsufficientCredits
    }

    return nil
}
```

**Принципы:**

- Repository **только** выполняет SQL запросы
- Repository **НЕ содержит** бизнес-логику
- Repository **возвращает** domain errors (не SQL ошибки)
- Repository **мапит** SQL rows → Go structs

#### 5. Middleware

**Задача:** Обработка запросов перед/после handlers

**JWT Authentication Middleware:**

```go
func JWTAuth(jwtManager *utils.JWTManager) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 1️⃣ Извлекаем токен из заголовка Authorization
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
                return
            }

            // Формат: "Bearer <token>"
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization format")
                return
            }

            tokenString := parts[1]

            // 2️⃣ Валидируем и парсим токен
            claims, err := jwtManager.ValidateToken(tokenString)
            if err != nil {
                respondWithError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
                return
            }

            // 3️⃣ Сохраняем данные пользователя в контексте
            ctx := r.Context()
            ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
            ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
            ctx = context.WithValue(ctx, IsPremiumKey, claims.IsPremium)

            // 4️⃣ Передаем управление следующему handler'у
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Получить User ID из контекста
func GetUserIDFromContext(r *http.Request) (uuid.UUID, bool) {
    userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
    return userID, ok
}

// Проверить Premium статус
func IsPremiumUser(r *http.Request) bool {
    isPremium, ok := r.Context().Value(IsPremiumKey).(bool)
    return ok && isPremium
}
```

**Как использовать в handler:**

```go
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    // Получаем ID пользователя из контекста (установлен middleware)
    userID, ok := middleware.GetUserIDFromContext(r)
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context")
        return
    }

    // Используем userID для получения профиля
    profile, err := h.profileService.GetByUserID(r.Context(), userID)
    // ...
}
```

#### 6. Models (Data Structures)

**User Model:**

```go
type User struct {
    ID                  uuid.UUID `json:"id"`
    Email               string    `json:"email"`
    PasswordHash        string    `json:"-"` // ⚠️ НИКОГДА не отправляем клиенту!
    FullName            string    `json:"full_name"`
    FreeGenerationsLeft int       `json:"free_generations_left"`
    IsPremium           bool      `json:"is_premium"`
    CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
}
```

**Profile Model:**

```go
type Profile struct {
    ID          uuid.UUID `json:"id"`
    UserID      uuid.UUID `json:"user_id"`
    Phone       string    `json:"phone,omitempty"`        // опционально
    Location    string    `json:"location,omitempty"`
    LinkedInURL string    `json:"linkedin_url,omitempty"`
    GithubURL   string    `json:"github_url,omitempty"`
    WebsiteURL  string    `json:"website_url,omitempty"`
    Summary     string    `json:"summary,omitempty"`      // О себе
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Experience Model:**

```go
type Experience struct {
    ID           uuid.UUID `json:"id"`
    ProfileID    uuid.UUID `json:"profile_id"`
    Company      string    `json:"company"`
    Position     string    `json:"position"`
    StartDate    Date      `json:"start_date"`
    EndDate      NullDate  `json:"end_date,omitempty"`    // null если текущее место
    IsCurrent    bool      `json:"is_current"`            // работаю сейчас
    Description  string    `json:"description,omitempty"`
    Achievements []string  `json:"achievements,omitempty"` // список достижений
    CreatedAt    time.Time `json:"created_at"`
}
```

**Document Model:**

```go
type Document struct {
    ID             uuid.UUID              `json:"id"`
    UserID         uuid.UUID              `json:"user_id"`
    Type           string                 `json:"type"` // "resume" или "cover_letter"
    Title          string                 `json:"title"`
    Content        map[string]interface{} `json:"content"` // JSONB - гибкая структура
    TemplateID     string                 `json:"template_id,omitempty"`
    JobTitle       string                 `json:"job_title,omitempty"`
    CompanyName    string                 `json:"company_name,omitempty"`
    JobDescription string                 `json:"job_description,omitempty"`
    Status         string                 `json:"status"` // "draft" или "final"
    CreatedAt      time.Time              `json:"created_at"`
    UpdatedAt      time.Time              `json:"updated_at"`
}
```

#### 7. Utils (Utilities)

**JWT Manager:**

```go
type JWTManager struct {
    secretKey []byte
}

func NewJWTManager(secret string) *JWTManager {
    return &JWTManager{
        secretKey: []byte(secret),
    }
}

// Генерация Access Token (короткий срок жизни: 15 минут)
func (m *JWTManager) GenerateAccessToken(user *models.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id":    user.ID.String(),
        "email":      user.Email,
        "is_premium": user.IsPremium,
        "exp":        time.Now().Add(15 * time.Minute).Unix(),
        "iat":        time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(m.secretKey)
}

// Генерация Refresh Token (долгий срок жизни: 7 дней)
func (m *JWTManager) GenerateRefreshToken(user *models.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID.String(),
        "type":    "refresh",
        "exp":     time.Now().Add(168 * time.Hour).Unix(), // 7 дней
        "iat":     time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(m.secretKey)
}

// Валидация токена
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return m.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        // Извлекаем данные из claims
        userID, _ := uuid.Parse(claims["user_id"].(string))
        email := claims["email"].(string)
        isPremium := claims["is_premium"].(bool)

        return &Claims{
            UserID:    userID,
            Email:     email,
            IsPremium: isPremium,
        }, nil
    }

    return nil, fmt.Errorf("invalid token")
}
```

**Password Hashing:**

```go
// Хеширование пароля (bcrypt)
func HashPassword(password string) (string, error) {
    // bcrypt автоматически добавляет salt
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

// Проверка пароля
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

**Почему bcrypt:**

- Автоматический salt (защита от rainbow tables)
- Медленный алгоритм (защита от brute-force)
- Настраиваемая сложность (cost factor)

### API Endpoints

#### Публичные (без авторизации)

```
GET  /api/v1/health               # Health check
POST /api/v1/auth/register        # Регистрация
POST /api/v1/auth/login           # Вход
POST /api/v1/auth/refresh         # Обновление токена
```

#### Защищенные (требуют JWT)

**Profile:**
```
GET    /api/v1/profile                    # Получить профиль
PUT    /api/v1/profile                    # Обновить профиль
POST   /api/v1/profile/experience         # Добавить опыт
PUT    /api/v1/profile/experience/{id}    # Обновить опыт
DELETE /api/v1/profile/experience/{id}    # Удалить опыт
POST   /api/v1/profile/education          # Добавить образование
PUT    /api/v1/profile/education/{id}     # Обновить образование
DELETE /api/v1/profile/education/{id}     # Удалить образование
POST   /api/v1/profile/skills             # Добавить навык
DELETE /api/v1/profile/skills/{id}        # Удалить навык
```

**Documents:**
```
POST   /api/v1/generate/resume            # Генерация резюме
POST   /api/v1/generate/cover-letter      # Генерация cover letter
GET    /api/v1/documents                  # Список документов
GET    /api/v1/documents/{id}             # Получить документ
PUT    /api/v1/documents/{id}             # Обновить документ
DELETE /api/v1/documents/{id}             # Удалить документ
GET    /api/v1/documents/{id}/pdf         # Экспорт в PDF
```

### Как запустить Backend

```bash
# 1. Переходим в папку backend
cd backend

# 2. Копируем .env.example в .env
cp .env.example .env

# 3. Редактируем .env (добавляем OpenAI API key)
nano .env

# 4. Запускаем PostgreSQL через Docker
make docker-up

# 5. Применяем миграции (создаем таблицы)
make migrate-up

# 6. Устанавливаем зависимости Go
make install

# 7. Запускаем сервер
make run

# Сервер запущен на http://localhost:8080
```

**Доступные make команды:**

```bash
make help          # Показать все команды
make run           # Запустить сервер
make build         # Собрать бинарник
make test          # Запустить тесты
make docker-up     # Запустить PostgreSQL
make docker-down   # Остановить контейнеры
make migrate-up    # Применить миграции
make db-shell      # Открыть psql shell
make lint          # Проверка кода
make fmt           # Форматирование кода
make dev           # Запустить все (БД + сервер)
```

---

## Frontend (React + TypeScript)

### Структура проекта

```
frontend/
├── public/
│   └── vite.svg
├── src/
│   ├── api/
│   │   ├── client.ts           # 🌐 Axios instance с interceptors
│   │   └── endpoints/
│   │       ├── auth.ts         # API методы для auth
│   │       ├── profile.ts      # API методы для profile
│   │       └── documents.ts    # API методы для documents
│   ├── app/
│   │   ├── App.tsx             # 🚪 Главный компонент
│   │   └── router.tsx          # 🛣️  Роутинг
│   ├── features/
│   │   ├── auth/
│   │   │   ├── components/
│   │   │   │   ├── LoginForm.tsx
│   │   │   │   └── RegisterForm.tsx
│   │   │   └── store/
│   │   │       └── authStore.ts     # 📦 Zustand store
│   │   ├── profile/
│   │   │   ├── components/
│   │   │   │   ├── ProfileForm.tsx
│   │   │   │   ├── ExperienceCard.tsx
│   │   │   │   └── SkillsSection.tsx
│   │   │   └── store/
│   │   │       └── profileStore.ts
│   │   └── documents/
│   │       ├── components/
│   │       │   ├── GenerateForm.tsx
│   │       │   └── DocumentCard.tsx
│   │       └── store/
│   │           └── documentsStore.ts
│   ├── pages/
│   │   ├── auth/
│   │   │   ├── LoginPage.tsx
│   │   │   └── RegisterPage.tsx
│   │   ├── Dashboard.tsx
│   │   ├── ProfilePage.tsx
│   │   ├── GeneratePage.tsx
│   │   └── HistoryPage.tsx
│   ├── shared/
│   │   ├── components/
│   │   │   ├── ui/              # UI компоненты
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Input.tsx
│   │   │   │   ├── Card.tsx
│   │   │   │   └── Modal.tsx
│   │   │   ├── Sidebar.tsx      # Боковое меню
│   │   │   ├── Header.tsx       # Шапка
│   │   │   └── ProtectedRoute.tsx # Защита роутов
│   │   ├── hooks/
│   │   │   └── useAuth.ts       # Custom hook для auth
│   │   ├── types/
│   │   │   └── api.types.ts     # TypeScript типы
│   │   └── utils/
│   │       └── format.ts        # Утилиты
│   ├── styles/
│   │   └── index.css            # TailwindCSS
│   └── main.tsx                 # 🚪 Entry point
├── .env                         # 🔐 Переменные окружения
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── tailwind.config.js
```

### Подробный разбор компонентов

#### 1. API Client (api/client.ts)

**Axios instance с interceptors:**

```typescript
import axios, { AxiosInstance, AxiosError } from 'axios';

// Создаем Axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request Interceptor (добавляем JWT токен к каждому запросу)
apiClient.interceptors.request.use(
  (config) => {
    // Получаем токен из localStorage
    const token = localStorage.getItem('access_token');

    if (token) {
      // Добавляем токен в заголовок Authorization
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response Interceptor (обработка ошибок)
apiClient.interceptors.response.use(
  (response) => {
    // Успешный ответ - просто возвращаем
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config;

    // Если получили 401 (Unauthorized)
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Пытаемся обновить токен
        const refreshToken = localStorage.getItem('refresh_token');
        const response = await axios.post('/api/v1/auth/refresh', {
          refresh_token: refreshToken,
        });

        const { access_token } = response.data;

        // Сохраняем новый токен
        localStorage.setItem('access_token', access_token);

        // Повторяем оригинальный запрос с новым токеном
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return apiClient(originalRequest);
      } catch (refreshError) {
        // Если обновление токена не удалось - logout
        localStorage.clear();
        window.location.href = '/auth/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;
```

**Почему это важно:**

- Автоматически добавляет JWT токен ко всем запросам
- Автоматически обновляет токен при истечении
- Автоматически делает logout при ошибке авторизации

#### 2. API Endpoints (api/endpoints/auth.ts)

```typescript
import apiClient from '../client';
import { User, TokenResponse } from '@/shared/types/api.types';

// Типы для requests
interface RegisterRequest {
  email: string;
  password: string;
  full_name: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

// API методы
export const authAPI = {
  // Регистрация
  register: async (data: RegisterRequest): Promise<TokenResponse> => {
    const response = await apiClient.post<TokenResponse>('/auth/register', data);
    return response.data;
  },

  // Вход
  login: async (data: LoginRequest): Promise<TokenResponse> => {
    const response = await apiClient.post<TokenResponse>('/auth/login', data);
    return response.data;
  },

  // Обновление токена
  refresh: async (refreshToken: string): Promise<{ access_token: string }> => {
    const response = await apiClient.post('/auth/refresh', {
      refresh_token: refreshToken,
    });
    return response.data;
  },

  // Logout (очистка локального хранилища)
  logout: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  },
};
```

#### 3. State Management (Zustand Store)

**Auth Store (features/auth/store/authStore.ts):**

```typescript
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authAPI } from '@/api/endpoints/auth';
import type { User } from '@/shared/types/api.types';

interface AuthState {
  // State
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  // Actions
  register: (email: string, password: string, fullName: string) => Promise<void>;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      // Начальное состояние
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      // Регистрация
      register: async (email, password, fullName) => {
        set({ isLoading: true, error: null });

        try {
          const response = await authAPI.register({
            email,
            password,
            full_name: fullName,
          });

          // Сохраняем токены в localStorage
          localStorage.setItem('access_token', response.access_token);
          localStorage.setItem('refresh_token', response.refresh_token);

          // Обновляем state
          set({
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error: any) {
          set({
            error: error.response?.data?.error?.message || 'Registration failed',
            isLoading: false,
          });
          throw error;
        }
      },

      // Вход
      login: async (email, password) => {
        set({ isLoading: true, error: null });

        try {
          const response = await authAPI.login({ email, password });

          localStorage.setItem('access_token', response.access_token);
          localStorage.setItem('refresh_token', response.refresh_token);

          set({
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error: any) {
          set({
            error: error.response?.data?.error?.message || 'Login failed',
            isLoading: false,
          });
          throw error;
        }
      },

      // Выход
      logout: () => {
        authAPI.logout();
        set({
          user: null,
          isAuthenticated: false,
          error: null,
        });
      },

      // Очистка ошибки
      clearError: () => set({ error: null }),
    }),
    {
      name: 'auth-storage', // Ключ в localStorage
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
```

**Как использовать в компонентах:**

```typescript
function LoginPage() {
  const { login, isLoading, error } = useAuthStore();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await login(email, password);
      navigate('/dashboard'); // Редирект после успешного входа
    } catch (err) {
      // Ошибка обработана в store
      console.error('Login failed:', err);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {error && <div className="error">{error}</div>}
      <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
      <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button type="submit" disabled={isLoading}>
        {isLoading ? 'Loading...' : 'Login'}
      </button>
    </form>
  );
}
```

#### 4. Protected Routes

**Компонент для защиты роутов (shared/components/ProtectedRoute.tsx):**

```typescript
import { Navigate } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store/authStore';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated } = useAuthStore();

  // Если пользователь не авторизован - редирект на login
  if (!isAuthenticated) {
    return <Navigate to="/auth/login" replace />;
  }

  // Если авторизован - рендерим страницу
  return <>{children}</>;
}
```

**Использование в роутере (app/router.tsx):**

```typescript
import { createBrowserRouter } from 'react-router-dom';
import { ProtectedRoute } from '@/shared/components/ProtectedRoute';
import LoginPage from '@/pages/auth/LoginPage';
import Dashboard from '@/pages/Dashboard';
import ProfilePage from '@/pages/ProfilePage';

export const router = createBrowserRouter([
  // Публичные роуты
  {
    path: '/auth/login',
    element: <LoginPage />,
  },
  {
    path: '/auth/register',
    element: <RegisterPage />,
  },

  // Защищенные роуты
  {
    path: '/dashboard',
    element: (
      <ProtectedRoute>
        <Dashboard />
      </ProtectedRoute>
    ),
  },
  {
    path: '/profile',
    element: (
      <ProtectedRoute>
        <ProfilePage />
      </ProtectedRoute>
    ),
  },
]);
```

#### 5. UI Components

**Button (shared/components/ui/Button.tsx):**

```typescript
import React from 'react';
import { cn } from '@/shared/utils/cn'; // утилита для объединения классов

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
}

export function Button({
  children,
  variant = 'primary',
  size = 'md',
  isLoading = false,
  className,
  disabled,
  ...props
}: ButtonProps) {
  const baseClasses = 'rounded-lg font-medium transition-colors focus:outline-none focus:ring-2';

  const variantClasses = {
    primary: 'bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500',
    secondary: 'bg-gray-200 hover:bg-gray-300 text-gray-900 focus:ring-gray-500',
    danger: 'bg-red-600 hover:bg-red-700 text-white focus:ring-red-500',
  };

  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
  };

  return (
    <button
      className={cn(
        baseClasses,
        variantClasses[variant],
        sizeClasses[size],
        (disabled || isLoading) && 'opacity-50 cursor-not-allowed',
        className
      )}
      disabled={disabled || isLoading}
      {...props}
    >
      {isLoading ? (
        <span className="flex items-center gap-2">
          <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
            {/* Spinner SVG */}
          </svg>
          Loading...
        </span>
      ) : (
        children
      )}
    </button>
  );
}
```

**Input (shared/components/ui/Input.tsx):**

```typescript
import React from 'react';
import { cn } from '@/shared/utils/cn';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export function Input({
  label,
  error,
  className,
  ...props
}: InputProps) {
  return (
    <div className="flex flex-col gap-1">
      {label && (
        <label className="text-sm font-medium text-gray-700">
          {label}
        </label>
      )}
      <input
        className={cn(
          'px-4 py-2 border rounded-lg focus:outline-none focus:ring-2',
          error
            ? 'border-red-500 focus:ring-red-500'
            : 'border-gray-300 focus:ring-blue-500',
          className
        )}
        {...props}
      />
      {error && (
        <span className="text-sm text-red-600">{error}</span>
      )}
    </div>
  );
}
```

### Как запустить Frontend

```bash
# 1. Переходим в папку frontend
cd frontend

# 2. Устанавливаем зависимости
npm install

# 3. Проверяем .env файл
cat .env
# Должно быть: VITE_API_BASE_URL=http://localhost:8080/api/v1

# 4. Запускаем dev сервер
npm run dev

# Frontend запущен на http://localhost:5173
```

**Доступные npm команды:**

```bash
npm run dev        # Запустить dev сервер
npm run build      # Собрать production build
npm run preview    # Просмотр production build
npm run lint       # Проверка кода
npm run format     # Форматирование
```

---

## База данных

### Схема базы данных

```sql
-- Таблица пользователей
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    free_generations_left INTEGER DEFAULT 2,
    is_premium BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица профилей (1:1 с users)
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phone VARCHAR(50),
    location VARCHAR(255),
    linkedin_url VARCHAR(255),
    github_url VARCHAR(255),
    website_url VARCHAR(255),
    summary TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица опыта работы
CREATE TABLE experiences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    company VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    is_current BOOLEAN DEFAULT FALSE,
    description TEXT,
    achievements TEXT[], -- Массив строк в PostgreSQL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица образования
CREATE TABLE education (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    institution VARCHAR(255) NOT NULL,
    degree VARCHAR(255) NOT NULL,
    field_of_study VARCHAR(255),
    start_date DATE NOT NULL,
    end_date DATE,
    gpa DECIMAL(3, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица навыков
CREATE TABLE skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50), -- technical, soft, language
    proficiency_level VARCHAR(50), -- beginner, intermediate, advanced, expert
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица документов
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- resume, cover_letter
    title VARCHAR(255) NOT NULL,
    content JSONB NOT NULL, -- Гибкая JSON структура
    template_id VARCHAR(50),
    job_title VARCHAR(255),
    company_name VARCHAR(255),
    job_description TEXT,
    status VARCHAR(50) DEFAULT 'draft', -- draft, final
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица истории генерации
CREATE TABLE generation_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    prompt_tokens INTEGER NOT NULL,
    completion_tokens INTEGER NOT NULL,
    total_cost DECIMAL(10, 4) NOT NULL,
    generation_time_ms INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для производительности
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_experiences_profile_id ON experiences(profile_id);
CREATE INDEX idx_education_profile_id ON education(profile_id);
CREATE INDEX idx_skills_profile_id ON skills(profile_id);
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_generation_history_user_id ON generation_history(user_id);
CREATE INDEX idx_users_email ON users(email);
```

### Связи между таблицами

```
users (1) ──── (1) profiles
users (1) ──── (N) documents
users (1) ──── (N) generation_history
profiles (1) ──── (N) experiences
profiles (1) ──── (N) education
profiles (1) ──── (N) skills
documents (1) ──── (N) generation_history
```

### Ключевые особенности

1. **UUID как Primary Keys**
   - Уникальные ID генерируются базой данных
   - Безопасно для распределенных систем

2. **CASCADE DELETE**
   - При удалении пользователя удаляются все связанные данные
   - При удалении профиля удаляются experiences, education, skills

3. **JSONB для content**
   - Гибкая структура документов
   - Можно хранить разные форматы резюме
   - Поддержка индексов и запросов по JSON

4. **TEXT[] для achievements**
   - Массив достижений в PostgreSQL
   - Удобно для хранения списков

---

## Аутентификация и безопасность

### JWT Flow

```
1. Регистрация / Вход
   ┌─────────┐                    ┌─────────┐
   │ Client  │ ─── POST /login ──>│ Backend │
   └─────────┘                    └─────────┘
                                       │
                                       ▼
                                  Проверка пароля
                                       │
                                       ▼
                                  Генерация токенов
                                  - Access Token (15min)
                                  - Refresh Token (7 days)
   ┌─────────┐                    ┌─────────┐
   │ Client  │ <── TokenResponse ─│ Backend │
   └─────────┘                    └─────────┘
        │
        ▼
   Сохраняет в localStorage
   - access_token
   - refresh_token

2. Защищенный запрос
   ┌─────────┐                              ┌─────────┐
   │ Client  │ ─── GET /profile ──────────>│ Backend │
   └─────────┘  Header: "Bearer <token>"   └─────────┘
                                                 │
                                                 ▼
                                            JWT Middleware
                                            - Проверка подписи
                                            - Проверка срока
                                            - Извлечение user_id
                                                 │
                                                 ▼
                                            Handler получает userID
                                            из контекста

3. Обновление токена
   ┌─────────┐                           ┌─────────┐
   │ Client  │ ─── POST /refresh ──────>│ Backend │
   └─────────┘  Body: refresh_token     └─────────┘
                                              │
                                              ▼
                                         Валидация refresh token
                                              │
                                              ▼
                                         Генерация нового
                                         access token
   ┌─────────┐                           ┌─────────┐
   │ Client  │ <── {access_token} ──────│ Backend │
   └─────────┘                           └─────────┘
```

### Безопасность паролей

**НЕ ПРАВИЛЬНО:**
```go
// ❌ НИКОГДА НЕ ДЕЛАЙТЕ ТАК!
user.Password = req.Password // Хранение пароля в открытом виде
```

**ПРАВИЛЬНО:**
```go
// ✅ Всегда хешируем пароли
passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
user.PasswordHash = string(passwordHash)
```

**Проверка пароля:**
```go
// ✅ Сравниваем хеш с паролем
err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
if err != nil {
    return errors.New("invalid password")
}
```

### CORS Configuration

```go
corsHandler := cors.New(cors.Options{
    AllowedOrigins:   []string{"http://localhost:5173"}, // Frontend URL
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"*"},
    AllowCredentials: true, // Важно для cookies/auth
})
```

### Защита от атак

1. **SQL Injection** - используем prepared statements
2. **XSS** - фронтенд экранирует все пользовательские данные
3. **CSRF** - JWT токены в headers (не в cookies)
4. **Brute Force** - bcrypt медленный (защита от подбора паролей)
5. **Rate Limiting** - ограничение количества запросов (TODO)

---

## Интеграция с OpenAI

### OpenAI Service

```go
type OpenAIService struct {
    apiKey string
    client *http.Client
}

// Генерация резюме
func (s *OpenAIService) GenerateResume(profile *models.Profile, jobDescription string) (string, error) {
    // 1️⃣ Формируем промпт
    prompt := fmt.Sprintf(`
You are a professional resume writer. Generate a tailored resume based on the following information:

Job Description:
%s

Candidate Profile:
- Name: %s
- Summary: %s
- Experience: %s
- Education: %s
- Skills: %s

Please generate a professional resume in JSON format with the following structure:
{
  "summary": "...",
  "experience": [...],
  "education": [...],
  "skills": [...]
}
`, jobDescription, profile.FullName, profile.Summary, formatExperience(profile), formatEducation(profile), formatSkills(profile))

    // 2️⃣ Создаем запрос к OpenAI API
    requestBody := map[string]interface{}{
        "model": "gpt-4",
        "messages": []map[string]string{
            {
                "role":    "system",
                "content": "You are a professional resume writer.",
            },
            {
                "role":    "user",
                "content": prompt,
            },
        },
        "temperature": 0.7,
        "max_tokens":  2000,
    }

    body, _ := json.Marshal(requestBody)

    // 3️⃣ Отправляем запрос
    req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

    resp, err := s.client.Do(req)
    if err != nil {
        return "", fmt.Errorf("openai request failed: %w", err)
    }
    defer resp.Body.Close()

    // 4️⃣ Парсим ответ
    var openAIResponse OpenAIResponse
    if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }

    // 5️⃣ Возвращаем сгенерированный контент
    return openAIResponse.Choices[0].Message.Content, nil
}
```

### Prompt Engineering

**Эффективный промпт для генерации резюме:**

```
1. Контекст: "You are a professional resume writer"
2. Входные данные: Job Description + Profile
3. Инструкции: "Generate a tailored resume"
4. Формат вывода: JSON structure
5. Tone: Professional, ATS-friendly
```

### Отслеживание использования

```go
// Сохраняем историю генерации
history := &models.GenerationHistory{
    UserID:           userID,
    DocumentID:       documentID,
    PromptTokens:     openAIResponse.Usage.PromptTokens,
    CompletionTokens: openAIResponse.Usage.CompletionTokens,
    TotalCost:        calculateCost(openAIResponse.Usage),
    GenerationTimeMs: int(time.Since(startTime).Milliseconds()),
}

s.historyRepo.Create(ctx, history)
```

---

## Разработка и тестирование

### Unit Tests

**Пример теста для AuthService:**

```go
func TestAuthService_Register(t *testing.T) {
    // Arrange (подготовка)
    mockUserRepo := &MockUserRepository{}
    mockProfileRepo := &MockProfileRepository{}
    jwtManager := utils.NewJWTManager("test-secret")
    authService := service.NewAuthService(mockUserRepo, mockProfileRepo, jwtManager)

    req := models.RegisterRequest{
        Email:    "test@example.com",
        Password: "password123",
        FullName: "Test User",
    }

    // Мокируем методы репозитория
    mockUserRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, repository.ErrUserNotFound)
    mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    mockProfileRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    // Act (действие)
    user, tokens, err := authService.Register(context.Background(), req)

    // Assert (проверка)
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.NotNil(t, tokens)
    assert.Equal(t, req.Email, user.Email)
    assert.Equal(t, 2, user.FreeGenerationsLeft)

    mockUserRepo.AssertExpectations(t)
    mockProfileRepo.AssertExpectations(t)
}
```

### Integration Tests

**Тест для HTTP handler:**

```go
func TestAuthHandler_Register(t *testing.T) {
    // Подготовка тестовой БД
    db := setupTestDB(t)
    defer db.Close()

    // Создаем реальные зависимости
    userRepo := repository.NewUserRepository(db)
    profileRepo := repository.NewProfileRepository(db)
    jwtManager := utils.NewJWTManager("test-secret")
    authService := service.NewAuthService(userRepo, profileRepo, jwtManager)
    authHandler := handlers.NewAuthHandler(authService, jwtManager)

    // Создаем HTTP запрос
    reqBody := `{"email":"test@example.com","password":"password123","full_name":"Test User"}`
    req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")

    // Создаем ResponseRecorder для записи ответа
    rr := httptest.NewRecorder()

    // Вызываем handler
    authHandler.Register(rr, req)

    // Проверяем результат
    assert.Equal(t, http.StatusCreated, rr.Code)

    var response models.TokenResponse
    json.NewDecoder(rr.Body).Decode(&response)

    assert.NotEmpty(t, response.AccessToken)
    assert.NotEmpty(t, response.RefreshToken)
    assert.Equal(t, "test@example.com", response.User.Email)
}
```

### QA Testing

В проекте есть папка `backend/qa/` с готовыми тестами:

```bash
# Запуск QA тестов
cd backend/qa
go test -v ./...
```

---

## Deployment

### Docker Compose (для разработки)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: resumebuilder-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: resume_builder
    ports:
      - "5433:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### Production Deployment

**Backend:**
1. Сборка: `make build`
2. Deploy бинарника на сервер
3. Настроить systemd service
4. Использовать reverse proxy (nginx)
5. SSL сертификаты (Let's Encrypt)

**Frontend:**
1. Сборка: `npm run build`
2. Deploy в `dist/` на CDN или static hosting
3. Настроить CORS для production домена

**Database:**
1. Managed PostgreSQL (AWS RDS, DigitalOcean Databases)
2. Регулярные бэкапы
3. Репликация для высокой доступности

---

## Типичные задачи

### Добавление нового endpoint'а

**1. Создаем handler:**

```go
// internal/handlers/skill.go
func (h *ProfileHandler) AddSkill(w http.ResponseWriter, r *http.Request) {
    userID, _ := middleware.GetUserIDFromContext(r)

    var req struct {
        Name             string `json:"name"`
        Category         string `json:"category"`
        ProficiencyLevel string `json:"proficiency_level"`
    }

    json.NewDecoder(r.Body).Decode(&req)

    skill, err := h.profileService.AddSkill(r.Context(), userID, req)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, skill)
}
```

**2. Добавляем в service:**

```go
// internal/service/profile_service.go
func (s *ProfileService) AddSkill(ctx context.Context, userID uuid.UUID, req SkillRequest) (*models.Skill, error) {
    // Валидация
    if req.Name == "" {
        return nil, errors.New("skill name is required")
    }

    // Создаем skill
    skill := &models.Skill{
        ID:               uuid.New(),
        ProfileID:        profileID,
        Name:             req.Name,
        Category:         req.Category,
        ProficiencyLevel: req.ProficiencyLevel,
        CreatedAt:        time.Now(),
    }

    // Сохраняем в БД
    if err := s.profileRepo.CreateSkill(ctx, skill); err != nil {
        return nil, err
    }

    return skill, nil
}
```

**3. Добавляем в repository:**

```go
// internal/repository/profile_repo.go
func (r *ProfileRepository) CreateSkill(ctx context.Context, skill *models.Skill) error {
    query := `
        INSERT INTO skills (id, profile_id, name, category, proficiency_level, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err := r.db.ExecContext(ctx, query,
        skill.ID,
        skill.ProfileID,
        skill.Name,
        skill.Category,
        skill.ProficiencyLevel,
        skill.CreatedAt,
    )

    return err
}
```

**4. Регистрируем route:**

```go
// cmd/api/main.go
protected.HandleFunc("/profile/skills", profileHandler.AddSkill).Methods("POST")
```

**5. Добавляем на frontend:**

```typescript
// api/endpoints/profile.ts
export const profileAPI = {
  addSkill: async (data: { name: string; category: string; proficiency_level: string }) => {
    const response = await apiClient.post('/profile/skills', data);
    return response.data;
  },
};
```

### Debugging

**Backend:**
```bash
# Логи сервера
go run cmd/api/main.go

# Запуск с debug
dlv debug cmd/api/main.go
```

**Frontend:**
```bash
# Dev mode с hot reload
npm run dev

# Browser DevTools:
# - Network tab для API запросов
# - Console для логов
# - React DevTools для state
```

**Database:**
```bash
# Подключение к БД
make db-shell

# Просмотр таблиц
\dt

# Запрос
SELECT * FROM users LIMIT 10;
```

---

## Заключение

Этот проект демонстрирует:

1. **Clean Architecture** с разделением слоев
2. **RESTful API** best practices
3. **JWT Authentication** и безопасность
4. **React + TypeScript** современный frontend
5. **PostgreSQL** реляционная БД
6. **AI Integration** с OpenAI API

### Дальнейшее развитие:

- Добавить экспорт в PDF
- Реализовать Premium подписку
- Добавить шаблоны резюме
- Улучшить AI промпты
- Добавить тесты
- Настроить CI/CD

### Полезные ресурсы:

- [Go Documentation](https://go.dev/doc/)
- [React Documentation](https://react.dev/)
- [PostgreSQL Tutorial](https://www.postgresql.org/docs/)
- [JWT.io](https://jwt.io/)
- [OpenAI API Docs](https://platform.openai.com/docs/)

---

**Удачи в разработке!** 🚀
