package usecase

import (
	"cafe/domain"
	"context"
	"time"
)

type logoutUsecase struct {
	refreshTokenRepository domain.RefreshTokenRepository
	blacklistService       domain.BlacklistService
	contextTimeout         time.Duration
}

func NewLogoutUsecase(refreshTokenRepository domain.RefreshTokenRepository, blacklistService domain.BlacklistService, timeout time.Duration) domain.LogoutUsecase {
	return &logoutUsecase{
		refreshTokenRepository: refreshTokenRepository,
		blacklistService:       blacklistService,
		contextTimeout:         timeout,
	}
}

func (lu *logoutUsecase) Logout(c context.Context, tokenID string, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()

	// 1. Add to Blacklist
	_ = lu.blacklistService.Add(ctx, tokenID, expiry)

	// 2. Delete from DB
	return lu.refreshTokenRepository.DeleteByTokenID(ctx, tokenID)
}

func (lu *logoutUsecase) LogoutAll(c context.Context, userID string, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()

	// 1. Get all active tokens for this user
	tokens, err := lu.refreshTokenRepository.GetByUserID(ctx, userID)
	if err == nil && len(tokens) > 0 {
		tokenIDs := make([]string, len(tokens))
		for i, t := range tokens {
			tokenIDs[i] = t.TokenID
		}
		// 2. Batch Add to Blacklist
		_ = lu.blacklistService.AddBatch(ctx, tokenIDs, expiry)
	}

	// 3. Delete all sessions from DB
	return lu.refreshTokenRepository.DeleteByUserID(ctx, userID)
}
