package server

import (
	"context"
	"gin-boot-starter/config"
	"gin-boot-starter/core/otel"
	"gin-boot-starter/core/validators"
	"gin-boot-starter/infra/db"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogging() {
	cfg := config.GetConfig().Logging
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		log.Info().Msgf("configuration logging.level: invalid. %v", cfg.Level)
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	// logger
	var writer io.Writer
	if cfg.Format == "json" {
		writer = os.Stdout
	} else {
		writer = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	}
	log.Logger = zerolog.New(writer).With().Timestamp().Caller().Logger()
}

func registerCustomValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("gender", validators.GenderValidator)
		v.RegisterValidation("enum", validators.EnumValidator)
	}
}

// Start gin server at config port
func Start() {
	//
	config := config.GetConfig()
	if config.Application.Profile == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	//
	initLogging()
	//
	// init OTEL tracing
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	//
	otel.InitProviders(ctx)
	//
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatal().Msg("DB Connection Failed.")
		return
	}
	acl := InitializeAclService()
	err = acl.LoadPolicy()
	if err != nil {
		log.Fatal().Msg("Load ACL rules failed.")
		return
	}
	//
	r := NewRouter()
	//
	registerCustomValidators()
	// graceful restart or stop
	//
	srv := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: r.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("listen Error on : %s\n", config.Server.Port)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//
	<-quit
	//
	log.Info().Msg("Shutdown Server ...")
	otel.Shutdown(context.Background())
	db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server Shutdown:")
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Info().Msg("timeout of 5 seconds.")
	}
	log.Info().Msg("Server exiting")
}
