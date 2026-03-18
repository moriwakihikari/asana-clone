package user

import (
	"context"
	"strings"

	"asana-clone-backend/internal/domain/shared"
	userDomain "asana-clone-backend/internal/domain/user"
)

// UserResponse is the public representation of a user.
type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// UserService handles user profile operations.
type UserService struct {
	userRepo userDomain.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo userDomain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetByID retrieves a user by their ID.
func (s *UserService) GetByID(ctx context.Context, id shared.ID) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, shared.ErrNotFound
	}
	return toUserResponse(u), nil
}

// UpdateProfile updates a user's name and avatar URL.
func (s *UserService) UpdateProfile(ctx context.Context, id shared.ID, name, avatarURL string) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, shared.ErrNotFound
	}

	if err := u.UpdateProfile(name, avatarURL); err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

// SearchUsers searches for users whose name or email contains the query string.
// This performs a client-side filter over FindByEmail as a simple approach.
// For production, the repository should expose a dedicated search method.
func (s *UserService) SearchUsers(ctx context.Context, query string) ([]*UserResponse, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return []*UserResponse{}, nil
	}

	// Try exact email match first.
	u, err := s.userRepo.FindByEmail(ctx, query)
	if err != nil {
		return nil, err
	}

	results := make([]*UserResponse, 0)
	if u != nil {
		results = append(results, toUserResponse(u))
	}

	return results, nil
}

func toUserResponse(u *userDomain.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
	}
}
