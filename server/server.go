package server

import (
	"gin001/config"
	"gin001/infra/db"
	"log"

	"github.com/gin-gonic/gin"
)

// Start gin server at config port
func Start() {
	config := config.GetConfig()
	if config.GetString("app.profile") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	//
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatalln("DB Connect Failed.")
		return
	}
	//
	r := NewRouter()
	r.Run(":" + config.GetString("server.port"))
}
