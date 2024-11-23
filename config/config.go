package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

var config *viper.Viper
var basePath = "/Users/yamingdeng/goprojects/src/gin001"

// LoadConf is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func LoadConf(env string) error {
	var err error
	//
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("application")
	config.AddConfigPath(basePath + "/config")
	err = config.ReadInConfig()
	if err != nil {
		log.Println("error on parsing default configuration file", err)
		return err
	}

	envConfig := viper.New()
	envConfig.SetConfigType("yaml")
	envConfig.AddConfigPath(basePath + "/config")
	envConfig.SetConfigName("application." + env)
	err = envConfig.ReadInConfig()
	if err != nil {
		log.Println("error on parsing env configuration file")
		return err
	}

	config.MergeConfigMap(envConfig.AllSettings())
	return nil
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

// GetConfig return application config.
func GetConfig() *viper.Viper {
	return config
}
