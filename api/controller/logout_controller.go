package controller

import (
	"cafe/bootstrap"
	"cafe/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LogoutController struct {
	LogoutUsecase domain.LogoutUsecase
	Env           *bootstrap.Env
}

func (lc *LogoutController) Logout(c *gin.Context) {
	refreshTokenID, exists := c.Get("x-refresh-token-id")
	expiry := time.Duration(lc.Env.AccessTokenExpiryHour) * time.Hour

	if exists && refreshTokenID != "" {
		err := lc.LogoutUsecase.Logout(c, refreshTokenID.(string), expiry)
		if err != nil {
			lc.clearCookie(c)
			c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Successfully logged out"})
			return
		}

		lc.clearCookie(c)
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Successfully logged out"})
		return
	}

	userID, userExists := c.Get("x-user-id")
	if userExists && userID != "" {
		err := lc.LogoutUsecase.LogoutAll(c, userID.(string), expiry)
		if err != nil {
			lc.clearCookie(c)
			c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Successfully logged out"})
			return
		}

		lc.clearCookie(c)
		c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Successfully logged out"})
		return
	}

	c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid session info"})
}

func (lc *LogoutController) clearCookie(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", lc.Env.AppEnv == "production", true)
}
