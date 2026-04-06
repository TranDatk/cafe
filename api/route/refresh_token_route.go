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

func NewRefreshTokenRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	rtRepo := repository.NewRefreshTokenRepository(db, domain.TableUserRefreshToken)
	userRepo := repository.NewUserRepository(db, domain.TableUser)
	txRepo := repository.NewTransaction(db)
	rtUsecase := usecase.NewRefreshTokenUsecase(rtRepo, userRepo, txRepo, timeout)
	rtc := &controller.RefreshTokenController{
		RefreshTokenUsecase: rtUsecase,
		Env:                 env,
	}
	group.POST("/refresh", rtc.RefreshToken)
}
