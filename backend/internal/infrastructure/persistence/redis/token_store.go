package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenStore manages JWT refresh tokens and blacklisting in Redis.
type TokenStore struct {
	client *redis.Client
}

// NewTokenStore creates a new TokenStore.
func NewTokenStore(client *redis.Client) *TokenStore {
	return &TokenStore{client: client}
}

// StoreRefreshToken stores a refresh token JTI associated with a user, with an expiry.
func (s *TokenStore) StoreRefreshToken(ctx context.Context, userID, jti string, expiry time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", jti)
	return s.client.Set(ctx, key, userID, expiry).Err()
}

// IsTokenBlacklisted checks whether a token JTI has been blacklisted.
func (s *TokenStore) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", jti)
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// BlacklistToken adds a token JTI to the blacklist with an expiry matching the token's remaining lifetime.
func (s *TokenStore) BlacklistToken(ctx context.Context, jti string, expiry time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", jti)
	return s.client.Set(ctx, key, "1", expiry).Err()
}

// DeleteRefreshToken removes a stored refresh token by its JTI.
func (s *TokenStore) DeleteRefreshToken(ctx context.Context, jti string) error {
	key := fmt.Sprintf("refresh_token:%s", jti)
	return s.client.Del(ctx, key).Err()
}
