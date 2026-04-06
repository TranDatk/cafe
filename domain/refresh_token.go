package domain

import "context"

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RefreshTokenUsecase interface {
	GetUserByID(c context.Context, id string) (User, error)
	GetRefreshTokenByID(c context.Context, tokenID string) (UserRefreshToken, error)
	CreateAccessToken(user *User, refreshTokenID string, secret string, expiry int) (string, string, error)
	RotateRefreshToken(c context.Context, user *User, oldTokenID string, secret string, expiry int) (string, string, error)
}
