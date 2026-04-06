package domain

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	UserID         string `json:"user_id"`
	RefreshTokenID string `json:"refresh_token_id"`
	jwt.RegisteredClaims
}
