package server

import (
	"gin-boot-starter/api/controller"
	"gin-boot-starter/core/config"
	"gin-boot-starter/core/logging"
	"gin-boot-starter/core/middleware"

	// server as OAS
	_ "gin-boot-starter/docs"

	"github.com/gin-contrib/pprof"

	"github.com/rs/zerolog/log"
	swaggerFile "github.com/swaggo/files"      // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// NewRouter create API routers.
func NewRouter() *gin.Engine {
	//
	router := gin.New()
	router.SetTrustedProxies([]string{"::1"})
	//
	// router.Use(gin.Logger())
	// Add logger as a middleware
	router.Use(logging.LoggerWithOptions(&logging.Options{
		Logger:        &log.Logger,
		FieldsExclude: []string{logging.BodyFieldName, logging.PayloadFieldName},
	}))
	//
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.SecurityHandler())
	router.Use(otelgin.Middleware(config.GetConfig().OTEL.ServiceName))

	//
	health := controller.NewHealthController()
	router.GET("/health", health.Status)

	// pprof
	debugGroup := router.Group("/debug", func(c *gin.Context) {
		// if c.Request.Header.Get("Authorization") != "foobar" {
		// 	c.AbortWithStatus(http.StatusForbidden)
		// 	return
		// }
		c.Next()
	})
	pprof.RouteRegister(debugGroup, "pprof")

	userGroup := router.Group("users", middleware.AuthJwtHandler())
	{
		user := InitializeUserController()
		userGroup.GET("/v1/paging/:size/:page", user.GetUsers)
		userGroup.GET("/v1/:id", user.GetUser)
		userGroup.PUT("/v1/:id", user.UpdateUser)
		userGroup.GET("/v1/session", user.GetCurrentUser)
	}

	accountGroup := router.Group("accounts")
	{
		account := InitializeAccountController()
		accountGroup.POST("/v1/signup", account.SignUp)
		accountGroup.POST("/v1/signin", account.SignIn)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))

	return router
}
