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
