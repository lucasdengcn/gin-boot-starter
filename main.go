package main

import (
	"flag"
	"fmt"
	"gin001/config"
	"gin001/infra/db"
	"gin001/migrations"
	"gin001/server"
	"os"
)

func getEnvOrFlagValue(flagValue, envVarName string) string {
	if flagValue == "" {
		return os.Getenv(envVarName)
	}
	return flagValue
}

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
	flagEnv := flag.String("e", "dev", "")
	flagCfg := flag.String("cfg", "", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -cfg {path} -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	//
	var envName = getEnvOrFlagValue(*flagEnv, "APP_ENV")
	var cfgPath = getEnvOrFlagValue(*flagCfg, "APP_CFG")
	//
	fmt.Printf("running in %v, env: %v\n", cfgPath, envName)
	// load configuration
	config.LoadConf(cfgPath, envName)
	// connect db
	db.ConnectDB()
	// build up db schema
	migrations.Build()
	server.Start()
}
