package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"asana-clone-backend/config"
	authApp "asana-clone-backend/internal/application/auth"
	commentApp "asana-clone-backend/internal/application/comment"
	projectApp "asana-clone-backend/internal/application/project"
	sectionApp "asana-clone-backend/internal/application/section"
	taskApp "asana-clone-backend/internal/application/task"
	userApp "asana-clone-backend/internal/application/user"
	workspaceApp "asana-clone-backend/internal/application/workspace"
	authInfra "asana-clone-backend/internal/infrastructure/auth"
	"asana-clone-backend/internal/infrastructure/persistence/postgres"
	redisInfra "asana-clone-backend/internal/infrastructure/persistence/redis"
	httpServer "asana-clone-backend/internal/interfaces/http"
	"asana-clone-backend/internal/interfaces/http/handler"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration.
	cfg := config.Load()

	log.Printf("Starting Asana Clone API server...")

	// Connect to PostgreSQL with retry.
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 10; i++ {
		pool, err = postgres.NewPostgresPool(cfg.DB)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to PostgreSQL (attempt %d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to PostgreSQL after 10 attempts: %v", err)
	}
	defer pool.Close()
	log.Printf("Connected to PostgreSQL")

	// Run database migrations.
	if err := postgres.RunMigrations(cfg.DB.DSN(), "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Printf("Database migrations applied")

	// Connect to Redis.
	redisClient, err := redisInfra.NewRedisClient(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Printf("Connected to Redis")

	// Initialize infrastructure services.
	jwtService := authInfra.NewJWTService(cfg.JWT.Secret)
	_ = redisInfra.NewTokenStore(redisClient)

	// Initialize repositories.
	userRepo := postgres.NewUserRepository(pool)
	workspaceRepo := postgres.NewWorkspaceRepository(pool)
	projectRepo := postgres.NewProjectRepository(pool)
	sectionRepo := postgres.NewSectionRepository(pool)
	taskRepo := postgres.NewTaskRepository(pool)
	commentRepo := postgres.NewCommentRepository(pool)
	labelRepo := postgres.NewLabelRepository(pool)

	// Initialize application services.
	authService := authApp.NewAuthService(userRepo, jwtService)
	userService := userApp.NewUserService(userRepo)
	workspaceService := workspaceApp.NewWorkspaceService(workspaceRepo, userRepo)
	projectService := projectApp.NewProjectService(projectRepo, sectionRepo, workspaceRepo)
	sectionService := sectionApp.NewSectionService(sectionRepo)
	taskService := taskApp.NewTaskService(taskRepo, sectionRepo, labelRepo, userRepo)
	commentService := commentApp.NewCommentService(commentRepo, userRepo)
	labelService := handler.NewLabelService(labelRepo)

	// Initialize handlers.
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	workspaceHandler := handler.NewWorkspaceHandler(workspaceService)
	projectHandler := handler.NewProjectHandler(projectService)
	sectionHandler := handler.NewSectionHandler(sectionService)
	taskHandler := handler.NewTaskHandler(taskService)
	commentHandler := handler.NewCommentHandler(commentService)
	labelHandler := handler.NewLabelHandler(labelService)

	// Create HTTP server.
	router := httpServer.NewServer(
		cfg,
		jwtService,
		authHandler,
		userHandler,
		workspaceHandler,
		projectHandler,
		sectionHandler,
		taskHandler,
		commentHandler,
		labelHandler,
	)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine.
	go func() {
		log.Printf("HTTP server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Printf("Server stopped gracefully")
}
