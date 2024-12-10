package main

import (
	"flag"
	"fmt"
	"gin-boot-starter/core/config"
	"gin-boot-starter/core/migration"
	"gin-boot-starter/server"
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
	//
	flagEnv := flag.String("e", "dev", "active profile, eg. dev, sit, uat, staging, prod")
	flagCfg := flag.String("w", "", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -w {path} -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	//
	var envName = getEnvOrFlagValue(*flagEnv, "APP_ENV")
	var basePath = getEnvOrFlagValue(*flagCfg, "APP_BASE")
	//
	fmt.Printf("running in %v, env: %v\n", basePath, envName)
	// load configuration
	err := config.LoadConf(basePath, envName)
	if err != nil {
		os.Exit(1)
	}
	// build up db schema
	migration.Build()
	//
	server.Start()
}
