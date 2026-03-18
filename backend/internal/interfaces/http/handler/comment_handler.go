package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	commentApp "asana-clone-backend/internal/application/comment"
	httpErrors "asana-clone-backend/internal/interfaces/errors"
	"asana-clone-backend/internal/interfaces/http/middleware"

	"github.com/go-chi/chi/v5"
)

// CommentHandler handles comment-related endpoints.
type CommentHandler struct {
	commentService *commentApp.CommentService
}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler(commentService *commentApp.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

type createCommentRequest struct {
	Content string `json:"content"`
}

type updateCommentRequest struct {
	Content string `json:"content"`
}

// Create handles POST /tasks/{taskID}/comments.
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	comment, err := h.commentService.AddComment(r.Context(), taskID, userID, req.Content)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusCreated, comment)
}

// List handles GET /tasks/{taskID}/comments.
func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	taskID, err := parseUUIDParam(r, "taskID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid task ID format",
		})
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	comments, err := h.commentService.ListByTask(r.Context(), taskID, offset, limit)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, comments)
}

// Update handles PATCH /tasks/{taskID}/comments/{commentID}.
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	commentID, err := parseUUIDParam(r, "commentID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid comment ID format",
		})
		return
	}

	var req updateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	comment, err := h.commentService.EditComment(r.Context(), commentID, userID, req.Content)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, comment)
}

// Delete handles DELETE /tasks/{taskID}/comments/{commentID}.
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	commentID, err := parseUUIDParam(r, "commentID")
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid comment ID format",
		})
		return
	}

	if err := h.commentService.DeleteComment(r.Context(), commentID, userID); err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Mount registers comment routes on the given router.
// These are nested under /tasks/{taskID}/comments.
func (h *CommentHandler) Mount(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Patch("/{commentID}", h.Update)
	r.Delete("/{commentID}", h.Delete)
}
