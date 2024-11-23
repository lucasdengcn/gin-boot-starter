package server

import (
	"gin001/apis/controllers"

	"github.com/gin-gonic/gin"
)

// NewRouter create API routers.
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	health := controllers.NewHealthController()

	router.GET("/health", health.Status)

	// router.Use(middlewares.AuthMiddleware())

	userGroup := router.Group("users")
	{
		user := InitializeUserController()
		//
		userGroup.POST("/v1/signup", user.SignUp)
		userGroup.POST("/v1/signin", user.SignIn)
		userGroup.GET("/v1/:id", user.GetUser)
	}

	// v1 := router.Group("v1")
	// {
	// 	userGroup := v1.Group("user")
	// 	{
	// 		user := new(controllers.UserController)
	// 		userGroup.GET("/:id", user.Retrieve)
	// 	}
	// }

	return router
}
