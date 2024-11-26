package server

import (
	"gin001/config"
	"gin001/infra/db"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	r.Run(":" + config.Server.Port)
}
