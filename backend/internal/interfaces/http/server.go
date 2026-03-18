package http

import (
	"net/http"

	"asana-clone-backend/config"
	"asana-clone-backend/internal/infrastructure/auth"
	"asana-clone-backend/internal/interfaces/http/handler"
	"asana-clone-backend/internal/interfaces/http/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewServer creates and configures the chi router with all routes and middleware.
func NewServer(
	cfg *config.Config,
	jwtService *auth.JWTService,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	workspaceHandler *handler.WorkspaceHandler,
	projectHandler *handler.ProjectHandler,
	sectionHandler *handler.SectionHandler,
	taskHandler *handler.TaskHandler,
	commentHandler *handler.CommentHandler,
	labelHandler *handler.LabelHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware.
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))
	r.Use(chiMiddleware.Recoverer)

	// Health check.
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Public routes.
	r.Route("/api/v1/auth", func(r chi.Router) {
		authHandler.Mount(r)
	})

	// Protected routes.
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))

		r.Route("/api/v1/users", func(r chi.Router) {
			userHandler.Mount(r)
		})

		r.Route("/api/v1/workspaces", func(r chi.Router) {
			workspaceHandler.Mount(r)

			// Nested project routes under workspace.
			r.Route("/{workspaceID}/projects", func(r chi.Router) {
				projectHandler.Mount(r)
			})

			// Nested label routes under workspace.
			r.Route("/{workspaceID}/labels", func(r chi.Router) {
				labelHandler.Mount(r)
			})

			// My tasks under workspace.
			r.Route("/{workspaceID}", func(r chi.Router) {
				taskHandler.MountMyTasks(r)
			})
		})

		r.Route("/api/v1/projects/{projectID}/sections", func(r chi.Router) {
			sectionHandler.Mount(r)
		})

		r.Route("/api/v1/projects/{projectID}/tasks", func(r chi.Router) {
			taskHandler.Mount(r)
		})

		r.Route("/api/v1/tasks", func(r chi.Router) {
			taskHandler.MountDirect(r)
		})

		r.Route("/api/v1/tasks/{taskID}/comments", func(r chi.Router) {
			commentHandler.Mount(r)
		})
	})

	return r
}
