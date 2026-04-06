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

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, rdb *redis.Client, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.TableUser)
	rtr := repository.NewRefreshTokenRepository(db, domain.TableUserRefreshToken)
	br := repository.NewRedisBlacklistRepository(rdb, "blacklist:")
	tx := repository.NewTransaction(db)

	lc := &controller.LoginController{
		LoginUsecase: usecase.NewSigninUsecase(ur, rtr, br, tx, timeout),
		Env:          env,
	}
	group.POST("/login", lc.Login)
}
