package server

import (
	"gin001/apis/controllers"
	"gin001/core/logging"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// NewRouter create API routers.
func NewRouter() *gin.Engine {
	//
	// logger to use with gin
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(writer).With().Timestamp().Logger()
	log.Logger = zerolog.New(writer).With().Timestamp().Logger()
	//
	router := gin.New()
	router.SetTrustedProxies([]string{"::1"})
	//
	// router.Use(gin.Logger())
	// Add logger as a middleware
	router.Use(logging.LoggerWithOptions(&logging.Options{Name: "App", Logger: &logger}))
	//
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

	return router
}
