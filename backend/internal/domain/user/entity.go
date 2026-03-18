package user

import (
	"strings"
	"time"

	"asana-clone-backend/internal/domain/shared"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           shared.ID
	Name         string
	Email        string
	PasswordHash string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(name, email, password string) (*User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" {
		return nil, shared.NewValidationError("INVALID_NAME", "name must not be empty", "name")
	}
	if email == "" {
		return nil, shared.NewValidationError("INVALID_EMAIL", "email must not be empty", "email")
	}
	if len(password) < 8 {
		return nil, shared.NewValidationError("INVALID_PASSWORD", "password must be at least 8 characters", "password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, shared.NewDomainError("PASSWORD_HASH_FAILED", "failed to hash password")
	}

	now := time.Now()
	return &User{
		ID:           shared.NewID(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) UpdateProfile(name, avatarURL string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return shared.NewValidationError("INVALID_NAME", "name must not be empty", "name")
	}

	u.Name = name
	u.AvatarURL = strings.TrimSpace(avatarURL)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) VerifyPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return shared.NewDomainError("INVALID_CREDENTIALS", "invalid email or password")
	}
	return nil
}
