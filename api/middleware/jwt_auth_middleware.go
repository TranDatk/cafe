package middleware

import (
	"net/http"
	"strings"

	"cafe/domain"
	tokenutil "cafe/internal/token_util"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
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
			c.Set("x-user-id", claims.ID)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
		c.Abort()
	}
}
