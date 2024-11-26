package server

import (
	"context"
	"gin001/config"
	"gin001/core/validators"
	"gin001/infra/db"
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
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	format := config.GetConfig().Logging.Format
	// logger
	var writer io.Writer
	if format == "json" {
		writer = os.Stdout
	} else {
		writer = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	}
	log.Logger = zerolog.New(writer).Level(zerolog.DebugLevel).With().Timestamp().Caller().Logger()
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
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatal().Msg("DB Connection Failed.")
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
	log.Info().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
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
