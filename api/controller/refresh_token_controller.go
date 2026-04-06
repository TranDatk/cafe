package controller

import (
	"cafe/bootstrap"
	"cafe/domain"
	tokenutil "cafe/internal/token_util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RefreshTokenController struct {
	RefreshTokenUsecase domain.RefreshTokenUsecase
	Env                 *bootstrap.Env
}

func (rtc *RefreshTokenController) RefreshToken(c *gin.Context) {
	var request domain.RefreshTokenRequest

	err := c.ShouldBindJSON(&request)
	refreshToken := request.RefreshToken

	if err != nil || refreshToken == "" {
		if cookieToken, err := c.Cookie("refresh_token"); err == nil {
			refreshToken = cookieToken
		}
	}

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Refresh token is required"})
		return
	}

	claims, err := tokenutil.ParseToken(refreshToken, rtc.Env.RefreshTokenSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}

	dbToken, err := rtc.RefreshTokenUsecase.GetRefreshTokenByID(c, claims.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Refresh token is invalid"})
		return
	}

	user, err := rtc.RefreshTokenUsecase.GetUserByID(c, claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "User not found"})
		return
	}

	var newRefreshToken string
	refreshTokenID := claims.ID

	threshold := time.Hour * 24
	if time.Until(dbToken.ExpiresAt) < threshold {
		var err error
		newRefreshToken, refreshTokenID, err = rtc.RefreshTokenUsecase.RotateRefreshToken(c, &user, claims.ID, rtc.Env.RefreshTokenSecret, rtc.Env.RefreshTokenExpiryHour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to rotate refresh token"})
			return
		}
	}

	accessToken, _, err := rtc.RefreshTokenUsecase.CreateAccessToken(&user, refreshTokenID, rtc.Env.AccessTokenSecret, rtc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	})
}
