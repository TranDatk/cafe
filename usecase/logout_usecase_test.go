package usecase

import (
	"cafe/domain"
	"cafe/domain/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogout(t *testing.T) {
	mockRefreshTokenRepo := new(mocks.RefreshTokenRepository)
	mockBlacklistService := new(mocks.BlacklistService)
	timeout := time.Second * 2

	lu := NewLogoutUsecase(mockRefreshTokenRepo, mockBlacklistService, timeout)

	t.Run("success", func(t *testing.T) {
		tokenID := "test-token-id"
		expiry := time.Hour

		mockBlacklistService.On("Add", mock.Anything, tokenID, expiry).Return(nil).Once()
		mockRefreshTokenRepo.On("DeleteByTokenID", mock.Anything, tokenID).Return(nil).Once()

		err := lu.Logout(context.Background(), tokenID, expiry)

		assert.NoError(t, err)
		mockRefreshTokenRepo.AssertExpectations(t)
		mockBlacklistService.AssertExpectations(t)
	})
}

func TestLogoutAll(t *testing.T) {
	mockRefreshTokenRepo := new(mocks.RefreshTokenRepository)
	mockBlacklistService := new(mocks.BlacklistService)
	timeout := time.Second * 2

	lu := NewLogoutUsecase(mockRefreshTokenRepo, mockBlacklistService, timeout)

	t.Run("success", func(t *testing.T) {
		userID := "test-user-id"
		expiry := time.Hour
		tokens := []domain.UserRefreshToken{
			{TokenID: "t1"},
			{TokenID: "t2"},
		}

		mockRefreshTokenRepo.On("GetByUserID", mock.Anything, userID).Return(tokens, nil).Once()
		mockBlacklistService.On("AddBatch", mock.Anything, []string{"t1", "t2"}, expiry).Return(nil).Once()
		mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, userID).Return(nil).Once()

		err := lu.LogoutAll(context.Background(), userID, expiry)

		assert.NoError(t, err)
		mockRefreshTokenRepo.AssertExpectations(t)
		mockBlacklistService.AssertExpectations(t)
	})
}
