package route

import (
	"cafe/api/controller"
	"cafe/bootstrap"
	"cafe/domain"
	"cafe/repository"
	"cafe/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewLogoutRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, rdb *redis.Client, group *gin.RouterGroup) {
	rtr := repository.NewRefreshTokenRepository(db, domain.TableUserRefreshToken)
	br := repository.NewRedisBlacklistRepository(rdb, "blacklist:")

	lc := &controller.LogoutController{
		LogoutUsecase: usecase.NewLogoutUsecase(rtr, br, timeout),
		Env:           env,
	}
	group.POST("/logout", lc.Logout)
}
