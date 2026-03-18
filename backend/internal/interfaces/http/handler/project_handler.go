package handler

import (
	"encoding/json"
	"net/http"

	projectApp "asana-clone-backend/internal/application/project"
	"asana-clone-backend/internal/domain/shared"
	httpErrors "asana-clone-backend/internal/interfaces/errors"
	"asana-clone-backend/internal/interfaces/http/middleware"

	"github.com/go-chi/chi/v5"
)

// ProjectHandler handles project-related endpoints.
type ProjectHandler struct {
	projectService *projectApp.ProjectService
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(projectService *projectApp.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

type createProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	ViewType    string `json:"view_type"`
}

type updateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	ViewType    string `json:"view_type"`
}

// Create handles POST /workspaces/{workspaceID}/projects.
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	project, err := h.projectService.Create(r.Context(), wsID, userID, req.Name, req.Description, req.Color, req.ViewType)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusCreated, project)
}

// List handles GET /workspaces/{workspaceID}/projects.
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	projects, err := h.projectService.ListByWorkspace(r.Context(), wsID, userID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, projects)
}

// Get handles GET /workspaces/{workspaceID}/projects/{projectID}.
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	project, err := h.projectService.GetByID(r.Context(), projectID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, project)
}

// Update handles PATCH /workspaces/{workspaceID}/projects/{projectID}.
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	var req updateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	project, err := h.projectService.Update(r.Context(), projectID, req.Name, req.Description, req.Color, req.ViewType)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, project)
}

// Delete handles DELETE /workspaces/{workspaceID}/projects/{projectID}.
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	if err := h.projectService.Delete(r.Context(), projectID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Archive handles POST /workspaces/{workspaceID}/projects/{projectID}/archive.
func (h *ProjectHandler) Archive(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	if err := h.projectService.Archive(r.Context(), projectID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "archived"})
}

// Unarchive handles POST /workspaces/{workspaceID}/projects/{projectID}/unarchive.
func (h *ProjectHandler) Unarchive(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	if err := h.projectService.Unarchive(r.Context(), projectID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "unarchived"})
}

// Mount registers project routes on the given router.
// These are nested under /workspaces/{workspaceID}/projects.
func (h *ProjectHandler) Mount(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{projectID}", h.Get)
	r.Patch("/{projectID}", h.Update)
	r.Delete("/{projectID}", h.Delete)
	r.Post("/{projectID}/archive", h.Archive)
	r.Post("/{projectID}/unarchive", h.Unarchive)
}

// parseUUIDParam extracts and parses a UUID URL parameter.
func parseUUIDParam(r *http.Request, param string) (shared.ID, error) {
	return shared.ParseID(chi.URLParam(r, param))
}
