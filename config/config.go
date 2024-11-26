package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var appConfig *Configuration

type Configuration struct {
	Application *Application
	DataSource  *DataSource
	Server      *Server
	Logging     *Logging
}

type Application struct {
	Name        string
	Description string
	Profile     string
}

type DataSource struct {
	URL     string
	Driver  string
	PoolMax int
	PoolMin int
}

type Server struct {
	Port string
}

type Logging struct {
	Format string
	Output string
}

// LoadConf is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func LoadConf(basePath, env string) error {
	var err error
	var config *viper.Viper
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
	//
	_appConfig := &Configuration{
		Application: &Application{
			Name:        config.GetString("app.name"),
			Description: config.GetString("app.description"),
			Profile:     config.GetString("app.profile"),
		},
		DataSource: &DataSource{
			URL:     config.GetString("datasource.url"),
			Driver:  config.GetString("datasource.driver"),
			PoolMax: config.GetInt("datasource.pool.max"),
			PoolMin: config.GetInt("datasource.pool.min"),
		},
		Server: &Server{
			Port: config.GetString("server.port"),
		},
		Logging: &Logging{
			Format: config.GetString("logging.format"),
			Output: config.GetString("logging.output"),
		},
	}
	// keep in global
	appConfig = _appConfig
	log.Printf("Configuration load success. %v\n", appConfig.Application.Name)
	return nil
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetBasePath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for _, err := os.ReadFile(filepath.Join(dir, "go.mod")); err != nil && len(dir) > 1; {
		println(dir)
		dir = filepath.Dir(dir)
		_, err = os.ReadFile(filepath.Join(dir, "go.mod"))
	}
	if len(dir) < 2 {
		panic("No go.mod found")
	}
	return dir
}

// GetConfig return application config.
func GetConfig() *Configuration {
	return appConfig
}
