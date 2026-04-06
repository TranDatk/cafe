package usecase

import (
	"cafe/domain"
	"context"
	"time"

	tokenutil "cafe/internal/token_util"

	"golang.org/x/crypto/bcrypt"
)

type signinUsecase struct {
	userRepository         domain.UserRepository
	refreshTokenRepository domain.RefreshTokenRepository
	blacklistService       domain.BlacklistService
	transaction            domain.Transaction
	contextTimeout         time.Duration
}

func NewSigninUsecase(
	userRepository domain.UserRepository,
	refreshTokenRepository domain.RefreshTokenRepository,
	blacklistService domain.BlacklistService,
	transaction domain.Transaction,
	timeout time.Duration) domain.LoginUsecase {

	return &signinUsecase{
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
		blacklistService:       blacklistService,
		transaction:            transaction,
		contextTimeout:         timeout,
	}
}

func (su *signinUsecase) Login(c context.Context, loginRequest *domain.LoginRequest) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	user, err := su.userRepository.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		return domain.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (su *signinUsecase) CreateAccessToken(user *domain.User, refreshTokenID string, secret string, expiry int) (string, string, error) {
	return tokenutil.CreateAccessToken(user, refreshTokenID, secret, expiry)
}

func (su *signinUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (string, string, error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}

func (su *signinUsecase) SaveRefreshToken(c context.Context, refreshToken *domain.UserRefreshToken, sessionCount int, accessTokenExpiryHour int) error {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	return su.transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		// 1. Check current active sessions
		count, err := su.refreshTokenRepository.GetCountByUserID(txCtx, refreshToken.UserID)
		if err != nil {
			return err
		}

		// 2. If limit reached, calculate how many to evict to make room for 1 new session
		toEvict := count - sessionCount + 1
		if toEvict > 0 {
			oldestSessions, err := su.refreshTokenRepository.GetOldestByUserID(txCtx, refreshToken.UserID, toEvict)
			if err == nil && len(oldestSessions) > 0 {
				tokenIDs := make([]string, len(oldestSessions))
				for i, session := range oldestSessions {
					tokenIDs[i] = session.TokenID
				}

				_ = su.refreshTokenRepository.DeleteByTokenIDs(txCtx, tokenIDs)

				expiryDuration := time.Duration(accessTokenExpiryHour) * time.Hour
				_ = su.blacklistService.AddBatch(txCtx, tokenIDs, expiryDuration)
			}
		}

		return su.refreshTokenRepository.Create(txCtx, refreshToken)
	})
}
