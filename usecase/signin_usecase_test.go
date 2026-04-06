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
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_Success(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRefreshTokenRepo := new(mocks.RefreshTokenRepository)
	mockBlacklistService := new(mocks.BlacklistService)
	mockRowTx := new(mocks.Transaction)

	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := domain.User{
		ID:       "1",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	loginRequest := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	mockUserRepo.On("GetByEmail", mock.Anything, loginRequest.Email).Return(user, nil)

	usecase := NewSigninUsecase(mockUserRepo, mockRefreshTokenRepo, mockBlacklistService, mockRowTx, time.Second*2)

	resultUser, err := usecase.Login(context.Background(), loginRequest)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, resultUser.ID)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRefreshTokenRepo := new(mocks.RefreshTokenRepository)
	mockBlacklistService := new(mocks.BlacklistService)
	mockRowTx := new(mocks.Transaction)

	loginRequest := &domain.LoginRequest{
		Email:    "notfound@example.com",
		Password: "testpassword",
	}

	mockUserRepo.On("GetByEmail", mock.Anything, loginRequest.Email).Return(domain.User{}, errors.New("user not found"))

	usecase := NewSigninUsecase(mockUserRepo, mockRefreshTokenRepo, mockBlacklistService, mockRowTx, time.Second*2)

	_, err := usecase.Login(context.Background(), loginRequest)

	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestSaveRefreshToken_BelowLimit(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRefreshTokenRepo := new(mocks.RefreshTokenRepository)
	mockBlacklistService := new(mocks.BlacklistService)
	mockRowTx := new(mocks.Transaction)

	refreshToken := &domain.UserRefreshToken{
		UserID:    "1",
		TokenID:   "token-123",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	sessionCount := 5
	accessTokenExpiryHour := 2

	// Mock transaction wrapper
	mockRowTx.On("WithinTransaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		_ = fn(context.Background())
	})

	mockRefreshTokenRepo.On("GetCountByUserID", mock.Anything, refreshToken.UserID).Return(2, nil)
	mockRefreshTokenRepo.On("Create", mock.Anything, refreshToken).Return(nil)

	usecase := NewSigninUsecase(mockUserRepo, mockRefreshTokenRepo, mockBlacklistService, mockRowTx, time.Second*2)

	err := usecase.SaveRefreshToken(context.Background(), refreshToken, sessionCount, accessTokenExpiryHour)

	assert.NoError(t, err)
	mockRefreshTokenRepo.AssertExpectations(t)
	mockRowTx.AssertExpectations(t)
}

func TestSaveRefreshToken_LimitReached(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRefreshTokenRepo := new(mocks.RefreshTokenRepository)
	mockBlacklistService := new(mocks.BlacklistService)
	mockRowTx := new(mocks.Transaction)

	refreshToken := &domain.UserRefreshToken{
		UserID:    "1",
		TokenID:   "new-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	oldestToken := domain.UserRefreshToken{
		UserID:  "1",
		TokenID: "oldest-token",
	}

	sessionCount := 5
	accessTokenExpiryHour := 2

	// Mock transaction wrapper
	mockRowTx.On("WithinTransaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		_ = fn(context.Background())
	})

	// It has exactly 5 sessions or more, so eviction hits
	mockRefreshTokenRepo.On("GetCountByUserID", mock.Anything, refreshToken.UserID).Return(5, nil)
	mockRefreshTokenRepo.On("GetOldestByUserID", mock.Anything, refreshToken.UserID, 1).Return([]domain.UserRefreshToken{oldestToken}, nil)
	mockRefreshTokenRepo.On("DeleteByTokenIDs", mock.Anything, []string{oldestToken.TokenID}).Return(nil)
	mockBlacklistService.On("AddBatch", mock.Anything, []string{oldestToken.TokenID}, time.Duration(accessTokenExpiryHour)*time.Hour).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.Anything, refreshToken).Return(nil)

	usecase := NewSigninUsecase(mockUserRepo, mockRefreshTokenRepo, mockBlacklistService, mockRowTx, time.Second*2)

	err := usecase.SaveRefreshToken(context.Background(), refreshToken, sessionCount, accessTokenExpiryHour)

	assert.NoError(t, err)
	mockRefreshTokenRepo.AssertExpectations(t)
	mockBlacklistService.AssertExpectations(t)
	mockRowTx.AssertExpectations(t)
}

func TestSaveRefreshToken_MultipleEvictions(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRefreshTokenRepo := new(mocks.RefreshTokenRepository)
	mockBlacklistService := new(mocks.BlacklistService)
	mockRowTx := new(mocks.Transaction)

	refreshToken := &domain.UserRefreshToken{
		UserID:    "1",
		TokenID:   "new-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	oldestTokens := []domain.UserRefreshToken{
		{UserID: "1", TokenID: "token-1-old"},
		{UserID: "1", TokenID: "token-2-old"},
		{UserID: "1", TokenID: "token-3-old"},
		{UserID: "1", TokenID: "token-4-old"},
	}

	sessionCount := 5
	accessTokenExpiryHour := 2

	// Mock transaction wrapper
	mockRowTx.On("WithinTransaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		_ = fn(context.Background())
	})

	// 8 existing sessions, limit 5. Need to evict: 8 - 5 + 1 = 4
	mockRefreshTokenRepo.On("GetCountByUserID", mock.Anything, refreshToken.UserID).Return(8, nil)
	mockRefreshTokenRepo.On("GetOldestByUserID", mock.Anything, refreshToken.UserID, 4).Return(oldestTokens, nil)

	tokenIDs := []string{"token-1-old", "token-2-old", "token-3-old", "token-4-old"}
	mockRefreshTokenRepo.On("DeleteByTokenIDs", mock.Anything, tokenIDs).Return(nil)
	mockBlacklistService.On("AddBatch", mock.Anything, tokenIDs, time.Duration(accessTokenExpiryHour)*time.Hour).Return(nil)

	mockRefreshTokenRepo.On("Create", mock.Anything, refreshToken).Return(nil)

	usecase := NewSigninUsecase(mockUserRepo, mockRefreshTokenRepo, mockBlacklistService, mockRowTx, time.Second*2)

	err := usecase.SaveRefreshToken(context.Background(), refreshToken, sessionCount, accessTokenExpiryHour)

	assert.NoError(t, err)
	mockRefreshTokenRepo.AssertExpectations(t)
	mockBlacklistService.AssertExpectations(t)
	mockRowTx.AssertExpectations(t)
}
