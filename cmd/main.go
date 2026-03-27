package main

import (
	"net/http"
	"time"

	"cafe/api/route"
	"cafe/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.App()
	defer app.CloseDBConnection()

	timeout := time.Duration(app.Env.ContextTimeout) * time.Second

	router := gin.Default()

	route.Setup(app.Env, timeout, app.DB, router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong test with RDS connection!",
		})
	})
	router.Run()
}
