package middleware

import (
	"cafe/domain"
	tokenutil "cafe/internal/token_util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string, blacklistService domain.BlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			claims, err := tokenutil.ParseToken(authToken, secret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
				c.Abort()
				return
			}

			// Blacklist Check
			isBlacklisted, _ := blacklistService.IsBlacklisted(c, claims.RefreshTokenID)
			if isBlacklisted {
				c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
				c.Abort()
				return
			}

			c.Set("x-user-id", claims.UserID)
			c.Set("x-refresh-token-id", claims.RefreshTokenID)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
		c.Abort()
	}
}
