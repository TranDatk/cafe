package domain

import (
	"context"
	"time"
)

type LogoutUsecase interface {
	Logout(c context.Context, tokenID string, expiry time.Duration) error
	LogoutAll(c context.Context, userID string, expiry time.Duration) error
}
