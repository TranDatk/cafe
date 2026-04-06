package domain

import (
	"context"
)

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginUsecase interface {
	Login(c context.Context, loginRequest *LoginRequest) (User, error)
	CreateAccessToken(user *User, refreshTokenID string, secret string, expiry int) (accessToken string, tokenID string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, tokenID string, err error)
	SaveRefreshToken(c context.Context, refreshToken *UserRefreshToken, sessionCount int, accessTokenExpiryHour int) error
}
