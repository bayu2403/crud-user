package main

import (
	"crud/user/controllers"
	"crud/user/models"
	"net/http"

	"crud/user/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gin-swagger middleware
// swagger embed files

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "User API"
	docs.SwaggerInfo.Description = "This is a CRUD User."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	route := gin.Default()

	models.ConnectDatabase()

	v1 := route.Group("/v1")
	{
		v1.GET("/v1/ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		v1.GET("/users", controllers.FindUsers)
		v1.POST("/users", controllers.CreateUsers)
		v1.GET("/users/:id", controllers.FindUser)
		v1.PATCH("/users/:id", controllers.UpdateUser)
		v1.DELETE("/users/:id", controllers.DeleteUser)
	}

	// use ginSwagger middleware to serve the API docs
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}
}
