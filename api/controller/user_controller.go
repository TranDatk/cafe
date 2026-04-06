package controller

import (
	"cafe/bootstrap"
	"cafe/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase domain.UserUsecase
	Env         *bootstrap.Env
}

func (uc *UserController) Fetch(c *gin.Context) {
	var opts domain.UserFetchOptions
	err := c.ShouldBind(&opts)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	opts.Prepare()

	users, total, err := uc.UserUsecase.Fetch(c, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.PaginationResponse[domain.User]{
		Data: users,
		Pagination: domain.Pagination{
			TotalRecords: total,
			TotalPages:   (total + int64(opts.PageSize) - 1) / int64(opts.PageSize),
			Page:         opts.Page,
			PageSize:     opts.PageSize,
		},
	})
}
