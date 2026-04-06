package route

import (
	"cafe/api/controller"
	"cafe/bootstrap"
	"cafe/domain"
	"cafe/repository"
	"cafe/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewLogoutRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, br domain.BlacklistService, group *gin.RouterGroup) {
	rtr := repository.NewRefreshTokenRepository(db, domain.TableUserRefreshToken)

	lc := &controller.LogoutController{
		LogoutUsecase: usecase.NewLogoutUsecase(rtr, br, timeout),
		Env:           env,
	}
	group.POST("/logout", lc.Logout)
}
