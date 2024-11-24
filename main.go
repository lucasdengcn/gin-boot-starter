package main

import (
	"flag"
	"fmt"
	"gin001/config"
	"gin001/infra/db"
	"gin001/migrations"
	"gin001/server"
	"log"
	"os"
)

// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	environment := flag.String("e", "dev", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	// load configuration
	config.LoadConf(*environment)
	cfg := config.GetConfig()
	fmt.Println(cfg.GetString("app.name"))
	// connect db
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatal("DB Connect Failed.")
	}
	// build up db schema
	migrations.Build()
	server.Start()
}
