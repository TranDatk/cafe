package route

import (
	"time"

	"cafe/api/middleware"
	"cafe/bootstrap"
	"cafe/domain"
	"cafe/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, rdb *redis.Client, gin *gin.Engine) {
	// Global Rate Limit
	rateLimitPerSec := 5.0
	burstSize := 10
	publicRouter := gin.Group("")
	publicRouter.Use(middleware.RateLimitMiddleware(rateLimitPerSec, burstSize))

	// All Public APIs
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, rdb, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)

	// All Private APIs
	protectedRouter := gin.Group("")

	// Initialize BlacklistService for Middleware
	br := repository.NewRedisBlacklistRepository(rdb, domain.BlacklistKeyPrefix)

	// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret, br))

	NewLogoutRouter(env, timeout, db, rdb, protectedRouter)
}
