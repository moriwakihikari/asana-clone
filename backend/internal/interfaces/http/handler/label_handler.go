package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"asana-clone-backend/internal/domain/label"
	"asana-clone-backend/internal/domain/shared"
	httpErrors "asana-clone-backend/internal/interfaces/errors"

	"github.com/go-chi/chi/v5"
)

// LabelResponse is the public representation of a label.
type LabelResponse struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	CreatedAt   string `json:"created_at"`
}

// LabelService is an inline application service for labels.
type LabelService struct {
	labelRepo label.LabelRepository
}

// NewLabelService creates a new LabelService.
func NewLabelService(labelRepo label.LabelRepository) *LabelService {
	return &LabelService{labelRepo: labelRepo}
}

// Create creates a new label.
func (s *LabelService) Create(ctx context.Context, workspaceID shared.ID, name, color string) (*LabelResponse, error) {
	lbl, err := label.NewLabel(workspaceID, name, color)
	if err != nil {
		return nil, err
	}

	if err := s.labelRepo.Save(ctx, lbl); err != nil {
		return nil, err
	}

	return toLabelResponse(lbl), nil
}

// ListByWorkspace returns all labels in a workspace.
func (s *LabelService) ListByWorkspace(ctx context.Context, workspaceID shared.ID) ([]*LabelResponse, error) {
	labels, err := s.labelRepo.FindByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	results := make([]*LabelResponse, 0, len(labels))
	for _, lbl := range labels {
		results = append(results, toLabelResponse(lbl))
	}
	return results, nil
}

// GetByID retrieves a label by ID.
func (s *LabelService) GetByID(ctx context.Context, id shared.ID) (*LabelResponse, error) {
	lbl, err := s.labelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lbl == nil {
		return nil, shared.ErrNotFound
	}
	return toLabelResponse(lbl), nil
}

// Update updates a label's name and color.
func (s *LabelService) Update(ctx context.Context, id shared.ID, name, color string) (*LabelResponse, error) {
	lbl, err := s.labelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lbl == nil {
		return nil, shared.ErrNotFound
	}

	if name != "" {
		lbl.Name = name
	}
	if color != "" {
		lbl.Color = color
	}

	if err := s.labelRepo.Save(ctx, lbl); err != nil {
		return nil, err
	}

	return toLabelResponse(lbl), nil
}

// Delete removes a label by ID.
func (s *LabelService) Delete(ctx context.Context, id shared.ID) error {
	lbl, err := s.labelRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if lbl == nil {
		return shared.ErrNotFound
	}

	return s.labelRepo.Delete(ctx, id)
}

func toLabelResponse(lbl *label.Label) *LabelResponse {
	return &LabelResponse{
		ID:          lbl.ID.String(),
		WorkspaceID: lbl.WorkspaceID.String(),
		Name:        lbl.Name,
		Color:       lbl.Color,
		CreatedAt:   lbl.CreatedAt.Format(time.RFC3339),
	}
}

// LabelHandler handles label-related endpoints.
type LabelHandler struct {
	labelService *LabelService
}

// NewLabelHandler creates a new LabelHandler.
func NewLabelHandler(labelService *LabelService) *LabelHandler {
	return &LabelHandler{labelService: labelService}
}

type createLabelRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type updateLabelRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Create handles POST /workspaces/{workspaceID}/labels.
func (h *LabelHandler) Create(w http.ResponseWriter, r *http.Request) {
	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	var req createLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	lbl, err := h.labelService.Create(r.Context(), wsID, req.Name, req.Color)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusCreated, lbl)
}

// List handles GET /workspaces/{workspaceID}/labels.
func (h *LabelHandler) List(w http.ResponseWriter, r *http.Request) {
	wsID, err := parseUUIDParam(r, "workspaceID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid workspace ID format",
		})
		return
	}

	labels, err := h.labelService.ListByWorkspace(r.Context(), wsID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, labels)
}

// Get handles GET /workspaces/{workspaceID}/labels/{labelID}.
func (h *LabelHandler) Get(w http.ResponseWriter, r *http.Request) {
	labelID, err := parseUUIDParam(r, "labelID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid label ID format",
		})
		return
	}

	lbl, err := h.labelService.GetByID(r.Context(), labelID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, lbl)
}

// Update handles PATCH /workspaces/{workspaceID}/labels/{labelID}.
func (h *LabelHandler) Update(w http.ResponseWriter, r *http.Request) {
	labelID, err := parseUUIDParam(r, "labelID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid label ID format",
		})
		return
	}

	var req updateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	lbl, err := h.labelService.Update(r.Context(), labelID, req.Name, req.Color)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, lbl)
}

// Delete handles DELETE /workspaces/{workspaceID}/labels/{labelID}.
func (h *LabelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	labelID, err := parseUUIDParam(r, "labelID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid label ID format",
		})
		return
	}

	if err := h.labelService.Delete(r.Context(), labelID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Mount registers label routes on the given router.
// These are nested under /workspaces/{workspaceID}/labels.
func (h *LabelHandler) Mount(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{labelID}", h.Get)
	r.Patch("/{labelID}", h.Update)
	r.Delete("/{labelID}", h.Delete)
}
