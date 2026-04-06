package usecase

import (
	"cafe/domain"
	tokenutil "cafe/internal/token_util"
	"context"
	"time"
)

type refreshTokenUsecase struct {
	refreshTokenRepository domain.RefreshTokenRepository
	userRepository         domain.UserRepository
	transaction            domain.Transaction
	contextTimeout         time.Duration
}

func NewRefreshTokenUsecase(refreshTokenRepository domain.RefreshTokenRepository, userRepository domain.UserRepository, transaction domain.Transaction, timeout time.Duration) domain.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		refreshTokenRepository: refreshTokenRepository,
		userRepository:         userRepository,
		transaction:            transaction,
		contextTimeout:         timeout,
	}
}

func (rtu *refreshTokenUsecase) GetUserByID(c context.Context, id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, rtu.contextTimeout)
	defer cancel()
	return rtu.userRepository.GetByID(ctx, id)
}

func (rtu *refreshTokenUsecase) GetRefreshTokenByID(c context.Context, tokenID string) (domain.UserRefreshToken, error) {
	ctx, cancel := context.WithTimeout(c, rtu.contextTimeout)
	defer cancel()
	return rtu.refreshTokenRepository.GetByTokenID(ctx, tokenID)
}

func (rtu *refreshTokenUsecase) CreateAccessToken(user *domain.User, refreshTokenID string, secret string, expiry int) (string, string, error) {
	return tokenutil.CreateAccessToken(user, refreshTokenID, secret, expiry)
}

func (rtu *refreshTokenUsecase) RotateRefreshToken(c context.Context, user *domain.User, oldTokenID string, secret string, expiry int) (string, string, error) {
	type rotateResult struct {
		token   string
		tokenID string
	}

	res, err := domain.WithinTransactionResult(rtu.transaction, c, func(txCtx context.Context) (rotateResult, error) {
		err := rtu.refreshTokenRepository.DeleteByTokenID(txCtx, oldTokenID)
		if err != nil {
			return rotateResult{}, err
		}

		refreshToken, refreshTokenID, err := tokenutil.CreateRefreshToken(user, secret, expiry)
		if err != nil {
			return rotateResult{}, err
		}

		err = rtu.refreshTokenRepository.Create(txCtx, &domain.UserRefreshToken{
			UserID:    user.ID,
			TokenID:   refreshTokenID,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expiry)),
		})

		if err != nil {
			return rotateResult{}, err
		}

		return rotateResult{token: refreshToken, tokenID: refreshTokenID}, nil
	})

	return res.token, res.tokenID, err
}
