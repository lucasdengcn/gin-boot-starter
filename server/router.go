package server

import (
	"gin001/apis/controllers"
	"gin001/core/logging"
	"gin001/core/middlewares"

	// server as OAS
	_ "gin001/docs"

	"github.com/rs/zerolog/log"
	swaggerFile "github.com/swaggo/files"      // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/gin-gonic/gin"
)

// NewRouter create API routers.
func NewRouter() *gin.Engine {
	//
	router := gin.New()
	router.SetTrustedProxies([]string{"::1"})
	//
	// router.Use(gin.Logger())
	// Add logger as a middleware
	router.Use(logging.LoggerWithOptions(&logging.Options{Name: "App", Logger: &log.Logger}))
	//
	router.Use(middlewares.ErrorHandler())
	router.Use(gin.Recovery())
	// Setup Security Headers
	router.Use(middlewares.SecurityHandler())
	//
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
		userGroup.GET("/v1/paging/:size/:page", user.GetUsers)
		userGroup.GET("/v1/:id", user.GetUser)
		userGroup.PUT("/v1/:id", user.UpdateUser)
	}

	// v1 := router.Group("v1")
	// {
	// 	userGroup := v1.Group("user")
	// 	{
	// 		user := new(controllers.UserController)
	// 		userGroup.GET("/:id", user.Retrieve)
	// 	}
	// }

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))

	return router
}
