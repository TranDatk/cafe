package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type BlacklistService struct {
	mock.Mock
}

func (m *BlacklistService) Add(c context.Context, tokenID string, expiration time.Duration) error {
	args := m.Called(c, tokenID, expiration)
	return args.Error(0)
}

func (m *BlacklistService) AddBatch(c context.Context, tokenIDs []string, expiration time.Duration) error {
	args := m.Called(c, tokenIDs, expiration)
	return args.Error(0)
}

func (m *BlacklistService) IsBlacklisted(c context.Context, tokenID string) (bool, error) {
	args := m.Called(c, tokenID)
	return args.Bool(0), args.Error(1)
}
