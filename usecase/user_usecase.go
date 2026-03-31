package usecase

import (
	"cafe/domain"
	"context"
	"time"
)

type userUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu *userUsecase) Fetch(c context.Context, opts domain.UserFetchOptions) ([]domain.User, int64, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	users, total, err := uu.userRepository.Fetch(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
