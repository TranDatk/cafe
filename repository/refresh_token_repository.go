package repository

import (
	"context"

	"cafe/domain"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	database *gorm.DB
	table    string
}

func NewRefreshTokenRepository(db *gorm.DB, table string) domain.RefreshTokenRepository {
	return &refreshTokenRepository{
		database: db,
		table:    table,
	}
}

func (rt *refreshTokenRepository) Create(c context.Context, refreshToken *domain.UserRefreshToken) error {
	return rt.database.WithContext(c).Table(rt.table).Create(refreshToken).Error
}

func (rt *refreshTokenRepository) DeleteByTokenID(c context.Context, tokenID string) error {
	return rt.database.WithContext(c).Table(rt.table).Where("token_id = ?", tokenID).Delete(&domain.UserRefreshToken{}).Error
}

func (rt *refreshTokenRepository) DeleteByTokenIDs(c context.Context, tokenIDs []string) error {
	return rt.database.WithContext(c).Table(rt.table).Where("token_id IN ?", tokenIDs).Delete(&domain.UserRefreshToken{}).Error
}

func (rt *refreshTokenRepository) GetByTokenID(c context.Context, tokenID string) (domain.UserRefreshToken, error) {
	var refreshToken domain.UserRefreshToken
	err := rt.database.WithContext(c).Table(rt.table).Where("token_id = ?", tokenID).First(&refreshToken).Error
	return refreshToken, err
}

func (rt *refreshTokenRepository) GetByUserID(c context.Context, userID string) ([]domain.UserRefreshToken, error) {
	var refreshTokens []domain.UserRefreshToken
	err := rt.database.WithContext(c).Table(rt.table).Where("user_id = ?", userID).Find(&refreshTokens).Error
	return refreshTokens, err
}

func (rt *refreshTokenRepository) DeleteByUserID(c context.Context, userID string) error {
	return rt.database.WithContext(c).Table(rt.table).Where("user_id = ?", userID).Delete(&domain.UserRefreshToken{}).Error
}

func (rt *refreshTokenRepository) GetCountByUserID(c context.Context, userID string) (int, error) {
	var count int64
	err := rt.database.WithContext(c).Table(rt.table).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}

func (rt *refreshTokenRepository) GetOldestByUserID(c context.Context, userID string, limit int) ([]domain.UserRefreshToken, error) {
	var refreshTokens []domain.UserRefreshToken
	err := rt.database.WithContext(c).Table(rt.table).Where("user_id = ?", userID).Order("created_at ASC").Limit(limit).Find(&refreshTokens).Error
	return refreshTokens, err
}
