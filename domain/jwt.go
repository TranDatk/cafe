package domain

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
