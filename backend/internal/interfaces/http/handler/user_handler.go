package handler

import (
	"encoding/json"
	"net/http"

	userApp "asana-clone-backend/internal/application/user"
	"asana-clone-backend/internal/domain/shared"
	httpErrors "asana-clone-backend/internal/interfaces/errors"
	"asana-clone-backend/internal/interfaces/http/middleware"

	"github.com/go-chi/chi/v5"
)

// UserHandler handles user-related endpoints.
type UserHandler struct {
	userService *userApp.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *userApp.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type updateProfileRequest struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// GetMe handles GET /users/me.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, user)
}

// UpdateMe handles PATCH /users/me.
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	user, err := h.userService.UpdateProfile(r.Context(), userID, req.Name, req.AvatarURL)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, user)
}

// GetUser handles GET /users/{userID}.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")
	id, err := shared.ParseID(idStr)
	if err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_ID",
			Message: "invalid user ID format",
		})
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, user)
}

// SearchUsers handles GET /users/search?q=.
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	users, err := h.userService.SearchUsers(r.Context(), query)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, users)
}

// Mount registers user routes on the given router.
func (h *UserHandler) Mount(r chi.Router) {
	r.Get("/me", h.GetMe)
	r.Patch("/me", h.UpdateMe)
	r.Get("/search", h.SearchUsers)
	r.Get("/{userID}", h.GetUser)
}
