package handler

import (
	"encoding/json"
	"net/http"

	authApp "asana-clone-backend/internal/application/auth"
	httpErrors "asana-clone-backend/internal/interfaces/errors"

	"github.com/go-chi/chi/v5"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *authApp.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *authApp.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Tokens *authApp.TokenPair    `json:"tokens"`
	User   *authApp.UserResponse `json:"user"`
}

// RegisterHandler handles POST /auth/register.
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	tokens, user, err := h.authService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusCreated, authResponse{
		Tokens: tokens,
		User:   user,
	})
}

// LoginHandler handles POST /auth/login.
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErrors.RespondWithJSON(w, http.StatusBadRequest, httpErrors.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
		return
	}

	tokens, user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httpErrors.RespondWithError(w, err)
		return
	}

	httpErrors.RespondWithJSON(w, http.StatusOK, authResponse{
		Tokens: tokens,
		User:   user,
	})
}

// Mount registers auth routes on the given router.
func (h *AuthHandler) Mount(r chi.Router) {
	r.Post("/register", h.RegisterHandler)
	r.Post("/login", h.LoginHandler)
}
