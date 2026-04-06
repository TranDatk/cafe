package tokenutil

import (
	"fmt"
	"time"

	"cafe/domain"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateAccessToken(user *domain.User, refreshTokenID string, secret string, expiry int) (accessToken string, tokenID string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()

	tID, err := uuid.NewV7()
	if err != nil {
		return "", "", err
	}
	tokenID = tID.String()

	claims := &domain.JwtCustomClaims{
		ID:             tokenID,
		Name:           user.Name,
		UserID:         user.ID,
		RefreshTokenID: refreshTokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return t, tokenID, err
}

func CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, tokenID string, err error) {
	tID, err := uuid.NewV7()
	if err != nil {
		return "", "", err
	}
	tokenID = tID.String()

	claimsRefresh := &domain.JwtCustomClaims{
		ID:     tokenID,
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Add(time.Hour*time.Duration(expiry)).Unix(), 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return rt, tokenID, err
}

func ParseToken(requestToken string, secret string) (*domain.JwtCustomClaims, error) {
	claims := &domain.JwtCustomClaims{}

	token, err := jwt.ParseWithClaims(
		requestToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
