package domain

import (
	"context"
	"time"
)

const (
	TableUserRefreshToken = "user_refresh_tokens"
	BlacklistKeyPrefix    = "blacklist:"
)

type UserRefreshToken struct {
	UserID    string    `gorm:"primaryKey;column:user_id" json:"user_id"`
	TokenID   string    `gorm:"primaryKey;column:token_id" json:"token_id"`
	ExpiresAt time.Time `gorm:"column:expires_at" json:"expires_at"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

type RefreshTokenRepository interface {
	Create(c context.Context, refreshToken *UserRefreshToken) error
	DeleteByTokenID(c context.Context, tokenID string) error
	DeleteByTokenIDs(c context.Context, tokenIDs []string) error
	GetByTokenID(c context.Context, tokenID string) (UserRefreshToken, error)
	GetByUserID(c context.Context, userID string) ([]UserRefreshToken, error)
	GetCountByUserID(c context.Context, userID string) (int, error)
	GetOldestByUserID(c context.Context, userID string, limit int) ([]UserRefreshToken, error)
	DeleteByUserID(c context.Context, userID string) error
}

type BlacklistService interface {
	Add(c context.Context, tokenID string, expiration time.Duration) error
	AddBatch(c context.Context, tokenIDs []string, expiration time.Duration) error
	IsBlacklisted(c context.Context, tokenID string) (bool, error)
}
