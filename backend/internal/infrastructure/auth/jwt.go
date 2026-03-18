package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenTTL  = 24 * time.Hour     // Dev-friendly: 24 hours
	RefreshTokenTTL = 7 * 24 * time.Hour // 7 days

	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Claims represents the JWT claims for both access and refresh tokens.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	JTI    string    `json:"jti"`
	Type   string    `json:"type"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secret string
}

// NewJWTService creates a new JWTService with the given signing secret.
func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: secret}
}

// GenerateTokenPair creates both an access token and a refresh token for the given user.
func (s *JWTService) GenerateTokenPair(userID uuid.UUID) (accessToken, refreshToken string, err error) {
	now := time.Now()

	// Generate access token.
	accessJTI := uuid.New().String()
	accessClaims := Claims{
		UserID: userID,
		JTI:    accessJTI,
		Type:   TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
			ID:        accessJTI,
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token.
	refreshJTI := uuid.New().String()
	refreshClaims := Claims{
		UserID: userID,
		JTI:    refreshJTI,
		Type:   TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenTTL)),
			ID:        refreshJTI,
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateAccessToken parses and validates an access token string.
func (s *JWTService) ValidateAccessToken(tokenStr string) (*Claims, error) {
	claims, err := s.parseToken(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.Type != TokenTypeAccess {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", TokenTypeAccess, claims.Type)
	}
	return claims, nil
}

// ValidateRefreshToken parses and validates a refresh token string.
func (s *JWTService) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	claims, err := s.parseToken(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.Type != TokenTypeRefresh {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", TokenTypeRefresh, claims.Type)
	}
	return claims, nil
}

func (s *JWTService) parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
