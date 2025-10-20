package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/feijoa-master/ai-resume-builder/internal/config"
	"github.com/feijoa-master/ai-resume-builder/internal/database"
	"github.com/feijoa-master/ai-resume-builder/internal/handlers"
	"github.com/feijoa-master/ai-resume-builder/internal/middleware"
	"github.com/feijoa-master/ai-resume-builder/internal/repository"
	"github.com/feijoa-master/ai-resume-builder/internal/service"
	"github.com/feijoa-master/ai-resume-builder/internal/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	log.Printf("Attempting to connect to database: %s", cfg.Database.URL)
	db, err := database.Connect(database.Config{
		URL:             cfg.Database.URL,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	profileRepo := repository.NewProfileRepository(db.DB)
	documentRepo := repository.NewDocumentRepository(db.DB)

	// Initialize OpenAI service
	openaiService := service.NewOpenAIService(
		cfg.OpenAI.APIKey,
		cfg.OpenAI.Model,
		cfg.OpenAI.MaxTokens,
		cfg.OpenAI.Temperature,
	)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	profileService := service.NewProfileService(profileRepo)
	documentService := service.NewDocumentService(documentRepo, profileRepo, userRepo, openaiService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(profileService)
	documentHandler := handlers.NewDocumentHandler(documentService)

	// Create router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check (public)
	api.HandleFunc("/health", healthCheckHandler(db)).Methods("GET")

	// Auth routes (public)
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.JWTAuth(jwtManager))

	// User endpoints
	protected.HandleFunc("/user/me", getMeHandler(authService)).Methods("GET")

	// Profile endpoints
	protected.HandleFunc("/profile", profileHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", profileHandler.UpdateProfile).Methods("PUT")

	// Experience endpoints
	protected.HandleFunc("/profile/experience", profileHandler.CreateExperience).Methods("POST")
	protected.HandleFunc("/profile/experience", profileHandler.GetExperiences).Methods("GET")
	protected.HandleFunc("/profile/experience/{id}", profileHandler.UpdateExperience).Methods("PUT")
	protected.HandleFunc("/profile/experience/{id}", profileHandler.DeleteExperience).Methods("DELETE")

	// Education endpoints
	protected.HandleFunc("/profile/education", profileHandler.CreateEducation).Methods("POST")
	protected.HandleFunc("/profile/education", profileHandler.GetEducation).Methods("GET")
	protected.HandleFunc("/profile/education/{id}", profileHandler.UpdateEducation).Methods("PUT")
	protected.HandleFunc("/profile/education/{id}", profileHandler.DeleteEducation).Methods("DELETE")

	// Skills endpoints
	protected.HandleFunc("/profile/skills", profileHandler.CreateSkill).Methods("POST")
	protected.HandleFunc("/profile/skills", profileHandler.GetSkills).Methods("GET")
	protected.HandleFunc("/profile/skills/{id}", profileHandler.DeleteSkill).Methods("DELETE")

	// Document generation endpoints
	protected.HandleFunc("/generate/resume", documentHandler.GenerateResume).Methods("POST")
	protected.HandleFunc("/generate/cover-letter", documentHandler.GenerateCoverLetter).Methods("POST")

	// Document management endpoints
	protected.HandleFunc("/documents", documentHandler.GetDocuments).Methods("GET")
	protected.HandleFunc("/documents/{id}", documentHandler.GetDocument).Methods("GET")
	protected.HandleFunc("/documents/{id}", documentHandler.UpdateDocument).Methods("PUT")
	protected.HandleFunc("/documents/{id}", documentHandler.DeleteDocument).Methods("DELETE")

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   cfg.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      corsHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  time.Second * 60,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

// Placeholder handlers (we'll implement these next)
func healthCheckHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Health(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status": "unhealthy"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "database": "connected"}`))
	}
}

func getMeHandler(authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserIDFromContext(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}

		user, err := authService.GetUserByID(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to get user"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
