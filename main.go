package main

import (
	"crud/user/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	route := gin.Default()

	models.ConnectDatabase()

	route.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}
}
