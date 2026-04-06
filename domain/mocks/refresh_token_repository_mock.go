package mocks

import (
	"cafe/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type RefreshTokenRepository struct {
	mock.Mock
}

func (m *RefreshTokenRepository) Create(c context.Context, refreshToken *domain.UserRefreshToken) error {
	args := m.Called(c, refreshToken)
	return args.Error(0)
}

func (m *RefreshTokenRepository) DeleteByTokenID(c context.Context, tokenID string) error {
	args := m.Called(c, tokenID)
	return args.Error(0)
}

func (m *RefreshTokenRepository) DeleteByTokenIDs(c context.Context, tokenIDs []string) error {
	args := m.Called(c, tokenIDs)
	return args.Error(0)
}

func (m *RefreshTokenRepository) GetByTokenID(c context.Context, tokenID string) (domain.UserRefreshToken, error) {
	args := m.Called(c, tokenID)
	return args.Get(0).(domain.UserRefreshToken), args.Error(1)
}

func (m *RefreshTokenRepository) GetByUserID(c context.Context, userID string) ([]domain.UserRefreshToken, error) {
	args := m.Called(c, userID)
	return args.Get(0).([]domain.UserRefreshToken), args.Error(1)
}

func (m *RefreshTokenRepository) DeleteByUserID(c context.Context, userID string) error {
	args := m.Called(c, userID)
	return args.Error(0)
}

func (m *RefreshTokenRepository) GetCountByUserID(c context.Context, userID string) (int, error) {
	args := m.Called(c, userID)
	return args.Get(0).(int), args.Error(1)
}

func (m *RefreshTokenRepository) GetOldestByUserID(c context.Context, userID string, limit int) ([]domain.UserRefreshToken, error) {
	args := m.Called(c, userID, limit)
	return args.Get(0).([]domain.UserRefreshToken), args.Error(1)
}
