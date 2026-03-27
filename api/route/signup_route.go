package route

import (
	"time"

	"cafe/api/controller"
	"cafe/bootstrap"
	"cafe/domain"
	"cafe/repository"
	"cafe/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.TableUser)
	rr := repository.NewRoleRepository(db, domain.TableRole)
	tx := repository.NewTransaction(db)

	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, rr, tx, timeout),
		Env:           env,
	}
	group.POST("/signup", sc.Signup)
}
