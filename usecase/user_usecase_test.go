package usecase

import (
	"cafe/domain"
	"cafe/domain/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch_Success(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	timeout := time.Second * 2
	usecase := NewUserUsecase(mockUserRepo, timeout)

	ctx := context.Background()
	opts := domain.UserFetchOptions{
		Page:     1,
		PageSize: 10,
	}

	expectedUsers := []domain.User{
		{ID: "1", Name: "Test User 1", Email: "test1@example.com"},
		{ID: "2", Name: "Test User 2", Email: "test2@example.com"},
	}
	expectedTotal := int64(2)

	// Since NewUserUsecase uses context.WithTimeout(c, timeout),
	// we use mock.Anything to match the internal context.
	mockUserRepo.On("Fetch", mock.Anything, opts).Return(expectedUsers, expectedTotal, nil)

	users, total, err := usecase.Fetch(ctx, opts)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	assert.Equal(t, expectedTotal, total)
	mockUserRepo.AssertExpectations(t)
}

func TestFetch_Error(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	timeout := time.Second * 2
	usecase := NewUserUsecase(mockUserRepo, timeout)

	ctx := context.Background()
	opts := domain.UserFetchOptions{
		Page:     1,
		PageSize: 10,
	}

	expectedErr := errors.New("repository error")

	mockUserRepo.On("Fetch", mock.Anything, opts).Return([]domain.User{}, int64(0), expectedErr)

	users, total, err := usecase.Fetch(ctx, opts)

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, expectedErr, err)
	mockUserRepo.AssertExpectations(t)
}
