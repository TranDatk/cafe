package controller

import (
	"cafe/bootstrap"
	"cafe/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

func (lc *LoginController) Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := lc.LoginUsecase.Login(c, &request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid email or password"})
		return
	}

	refreshToken, refreshTokenID, err := lc.LoginUsecase.CreateRefreshToken(&user, lc.Env.RefreshTokenSecret, lc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, _, err := lc.LoginUsecase.CreateAccessToken(&user, refreshTokenID, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err = lc.LoginUsecase.SaveRefreshToken(c, &domain.UserRefreshToken{
		UserID:    user.ID,
		TokenID:   refreshTokenID,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(lc.Env.RefreshTokenExpiryHour)),
	}, lc.Env.SessionCount, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.SetCookie("refresh_token", refreshToken, lc.Env.RefreshTokenExpiryHour*3600, "/", "", lc.Env.AppEnv == "production", true)

	c.JSON(http.StatusOK, loginResponse)
}
