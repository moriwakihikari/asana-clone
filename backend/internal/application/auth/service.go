package auth

import (
	"context"
	"strings"

	"asana-clone-backend/internal/domain/shared"
	"asana-clone-backend/internal/domain/user"

	"github.com/google/uuid"
)

// TokenGenerator abstracts JWT token generation (implemented in infrastructure).
type TokenGenerator interface {
	GenerateTokenPair(userID uuid.UUID) (accessToken, refreshToken string, err error)
}

// TokenPair holds the access and refresh tokens returned on auth.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserResponse is the public representation of a user returned by auth endpoints.
type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// AuthService handles registration and login.
type AuthService struct {
	userRepo       user.UserRepository
	tokenGenerator TokenGenerator
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo user.UserRepository, tokenGenerator TokenGenerator) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		tokenGenerator: tokenGenerator,
	}
}

// Register creates a new user account and returns tokens.
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*TokenPair, *UserResponse, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	// Check if the email is already taken.
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if existing != nil {
		return nil, nil, shared.NewValidationError("EMAIL_TAKEN", "a user with this email already exists", "email")
	}

	// Create the domain entity (handles validation + password hashing).
	u, err := user.NewUser(name, email, password)
	if err != nil {
		return nil, nil, err
	}

	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, nil, err
	}

	// Generate JWT tokens.
	accessToken, refreshToken, err := s.tokenGenerator.GenerateTokenPair(u.ID)
	if err != nil {
		return nil, nil, shared.NewDomainError("TOKEN_GENERATION_FAILED", "failed to generate authentication tokens")
	}

	tokens := &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	userResp := toUserResponse(u)

	return tokens, userResp, nil
}

// Login verifies credentials and returns tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, *UserResponse, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if u == nil {
		return nil, nil, shared.NewDomainError("INVALID_CREDENTIALS", "invalid email or password")
	}

	if err := u.VerifyPassword(password); err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, err := s.tokenGenerator.GenerateTokenPair(u.ID)
	if err != nil {
		return nil, nil, shared.NewDomainError("TOKEN_GENERATION_FAILED", "failed to generate authentication tokens")
	}

	tokens := &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	userResp := toUserResponse(u)

	return tokens, userResp, nil
}

func toUserResponse(u *user.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
	}
}
