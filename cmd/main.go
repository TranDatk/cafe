package main

import (
	"net/http"

	"cafe/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.App()
	defer app.CloseDBConnection()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong test with RDS connection!",
		})
	})
	router.Run()
}
