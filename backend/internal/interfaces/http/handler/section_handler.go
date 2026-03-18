package handler

import (
	"encoding/json"
	"net/http"

	sectionApp "asana-clone-backend/internal/application/section"
	httpErrors "asana-clone-backend/internal/interfaces/errors"

	"github.com/go-chi/chi/v5"
)

// SectionHandler handles section-related endpoints.
type SectionHandler struct {
	sectionService *sectionApp.SectionService
}

// NewSectionHandler creates a new SectionHandler.
func NewSectionHandler(sectionService *sectionApp.SectionService) *SectionHandler {
	return &SectionHandler{sectionService: sectionService}
}

type createSectionRequest struct {
	Name string `json:"name"`
}

type renameSectionRequest struct {
	Name string `json:"name"`
}

type reorderSectionRequest struct {
	Position int `json:"position"`
}

// Create handles POST /projects/{projectID}/sections.
func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	var req createSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	section, err := h.sectionService.Create(r.Context(), projectID, req.Name)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusCreated, section)
}

// List handles GET /projects/{projectID}/sections.
func (h *SectionHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDParam(r, "projectID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid project ID format",
		})
		return
	}

	sections, err := h.sectionService.ListByProject(r.Context(), projectID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, sections)
}

// Update handles PATCH /projects/{projectID}/sections/{sectionID}.
func (h *SectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	sectionID, err := parseUUIDParam(r, "sectionID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid section ID format",
		})
		return
	}

	var req renameSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	section, err := h.sectionService.Rename(r.Context(), sectionID, req.Name)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, section)
}

// Delete handles DELETE /projects/{projectID}/sections/{sectionID}.
func (h *SectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sectionID, err := parseUUIDParam(r, "sectionID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid section ID format",
		})
		return
	}

	if err := h.sectionService.Delete(r.Context(), sectionID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Reorder handles POST /projects/{projectID}/sections/{sectionID}/reorder.
func (h *SectionHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	sectionID, err := parseUUIDParam(r, "sectionID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid section ID format",
		})
		return
	}

	var req reorderSectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	if err := h.sectionService.Reorder(r.Context(), sectionID, req.Position); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "reordered"})
}

// Mount registers section routes on the given router.
// These are nested under /projects/{projectID}/sections.
func (h *SectionHandler) Mount(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Patch("/{sectionID}", h.Update)
	r.Delete("/{sectionID}", h.Delete)
	r.Post("/{sectionID}/reorder", h.Reorder)
}
