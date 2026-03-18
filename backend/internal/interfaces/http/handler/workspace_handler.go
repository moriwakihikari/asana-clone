package handler

import (
	"encoding/json"
	"net/http"

	workspaceApp "asana-clone-backend/internal/application/workspace"
	httpErrors "asana-clone-backend/internal/interfaces/errors"
	"asana-clone-backend/internal/interfaces/http/middleware"

	"github.com/go-chi/chi/v5"
)

// WorkspaceHandler handles workspace-related endpoints.
type WorkspaceHandler struct {
	workspaceService *workspaceApp.WorkspaceService
}

// NewWorkspaceHandler creates a new WorkspaceHandler.
func NewWorkspaceHandler(workspaceService *workspaceApp.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceService: workspaceService}
}

type createWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type addMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Create handles POST /workspaces.
func (h *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req createWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	ws, err := h.workspaceService.Create(r.Context(), userID, req.Name, req.Description)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusCreated, ws)
}

// List handles GET /workspaces.
func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	workspaces, err := h.workspaceService.ListByUser(r.Context(), userID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, workspaces)
}

// Get handles GET /workspaces/{workspaceID}.
func (h *WorkspaceHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	ws, err := h.workspaceService.GetByID(r.Context(), wsID, userID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, ws)
}

// Update handles PATCH /workspaces/{workspaceID}.
func (h *WorkspaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	var req updateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	if err := h.workspaceService.Update(r.Context(), wsID, userID, req.Name, req.Description); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// Delete handles DELETE /workspaces/{workspaceID}.
func (h *WorkspaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	if err := h.workspaceService.Delete(r.Context(), wsID, userID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListMembers handles GET /workspaces/{workspaceID}/members.
func (h *WorkspaceHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	ws, err := h.workspaceService.GetByID(r.Context(), wsID, userID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, ws.Members)
}

// AddMember handles POST /workspaces/{workspaceID}/members.
func (h *WorkspaceHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	var req addMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	if err := h.workspaceService.AddMember(r.Context(), wsID, req.Email, req.Role, userID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusCreated, map[string]string{"status": "member_added"})
}

// RemoveMember handles DELETE /workspaces/{workspaceID}/members/{userID}.
func (h *WorkspaceHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetUserID(r.Context())

	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	memberID, err := parseUUIDParam(r, "userID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid user ID format",
		})
		return
	}

	if err := h.workspaceService.RemoveMember(r.Context(), wsID, memberID, callerID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Mount registers workspace routes on the given router.
func (h *WorkspaceHandler) Mount(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{workspaceID}", h.Get)
	r.Patch("/{workspaceID}", h.Update)
	r.Delete("/{workspaceID}", h.Delete)
	r.Get("/{workspaceID}/members", h.ListMembers)
	r.Post("/{workspaceID}/members", h.AddMember)
	r.Delete("/{workspaceID}/members/{userID}", h.RemoveMember)
}
