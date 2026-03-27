package domain

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name   string `json:"name"`
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}
