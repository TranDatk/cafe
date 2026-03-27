package usecase

import (
	"cafe/domain"
	"context"
	"time"

	tokenutil "cafe/internal/token_util"
)

type signupUsecase struct {
	userRepository domain.UserRepository
	roleRepository domain.RoleRepository
	transaction    domain.Transaction
	contextTimeout time.Duration
}

func NewSignupUsecase(userRepository domain.UserRepository, roleRepository domain.RoleRepository, transaction domain.Transaction, timeout time.Duration) domain.SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
		roleRepository: roleRepository,
		transaction:    transaction,
		contextTimeout: timeout,
	}
}

func (su *signupUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	return su.userRepository.GetByEmail(ctx, email)
}

func (su *signupUsecase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	return su.transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		// 1. Create user
		if err := su.userRepository.Create(txCtx, user); err != nil {
			return err
		}

		// 2. Get default role ("user")
		role, err := su.roleRepository.GetByName(txCtx, domain.UserRole)
		if err != nil {
			return err
		}

		// 3. Assign role to user
		return su.userRepository.AssignRole(txCtx, user, &role)
	})
}

func (su *signupUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (string, error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (su *signupUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (string, error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}
