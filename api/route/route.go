package route

import (
	"time"

	"cafe/api/middleware"
	"cafe/bootstrap"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, gin *gin.Engine) {
	// Global Rate Limit
	rateLimitPerSec := 5.0
	burstSize := 10
	publicRouter := gin.Group("")
	publicRouter.Use(middleware.RateLimitMiddleware(rateLimitPerSec, burstSize))

	// All Public APIs
	NewSignupRouter(env, timeout, db, publicRouter)
	// NewLoginRouter(env, timeout, db, publicRouter)
	// NewRefreshTokenRouter(env, timeout, db, publicRouter)

	// protectedRouter := gin.Group("")
	// Middleware to verify AccessToken
	// protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// All Private APIs
	// NewProfileRouter(env, timeout, db, protectedRouter)
	// NewTaskRouter(env, timeout, db, protectedRouter)
}
