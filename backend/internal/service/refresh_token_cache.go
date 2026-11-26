package service

import (
	"context"
	"time"
)

// RefreshTokenData stores refresh-token metadata in cache.
type RefreshTokenData struct {
	UserID       int64     `json:"user_id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	TokenVersion int64     `json:"token_version"`
	FamilyID     string    `json:"family_id"`
	ParentHash   string    `json:"parent_hash,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshTokenCache abstracts refresh-token persistence (typically Redis).
type RefreshTokenCache interface {
	StoreRefreshToken(ctx context.Context, tokenHash string, data *RefreshTokenData, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenData, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	AddToUserTokenSet(ctx context.Context, userID int64, tokenHash string, ttl time.Duration) error
	AddToFamilyTokenSet(ctx context.Context, familyID, tokenHash string, ttl time.Duration) error
	DeleteTokenFamily(ctx context.Context, familyID string) error
	DeleteUserRefreshTokens(ctx context.Context, userID int64) error
}
