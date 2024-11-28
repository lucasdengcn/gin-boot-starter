package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

var appConfig *Configuration

type Configuration struct {
	Application *Application
	DataSource  *DataSource
	Server      *Server
	Logging     *Logging
	OTEL        *OTEL
	Security    *Security
}

type Application struct {
	Name        string
	Description string
	Profile     string
	CfgPath     string
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
	Level  string
	Format string
	Output string
}

type OTEL struct {
	ServiceName  string
	ServiceVer   string
	Insecure     string
	OTLPEndpoint string
	Logging      bool
	Tracer       bool
	Metric       bool
}

type JWT struct {
	/*
			private_key: "your-secret-key"
		    public_key: "your-secret-key"
		    expiration_time: 10 // minutes
		    refresh_token_expiration_time: 72 // hours
		    token_blacklist_enabled: true
		    token_blacklist_ttl: 30 // minutes
	*/
	PrivateKey                      string
	PublicKey                       string
	ExpirationTimeMinutes           int
	RefreshTokenExpirationTimeHours int
	TokenBlackListEnabled           bool
	TokenBlackListTTL               int
}

type Security struct {
	JWT *JWT
}

func value(v *viper.Viper, cfgKey, envKey, defaultValue string) string {
	cfgValue := os.Getenv(envKey)
	if cfgValue != "" {
		return cfgValue
	}
	cfgValue = v.GetString(cfgKey)
	if cfgValue != "" {
		return cfgValue
	}
	return defaultValue
}

func intValue(v *viper.Viper, cfgKey, envKey string, defaultValue int) int {
	cfgValue := value(v, cfgKey, envKey, "")
	if cfgValue == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(cfgValue)
	if err != nil {
		log.Printf("parse key %v, %v to int error.", cfgKey, envKey)
		return defaultValue
	}
	return val
}

// LoadConf is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func LoadConf(cfgPath, env string) error {
	var err error
	var config *viper.Viper
	//
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("application")
	config.AddConfigPath(cfgPath)
	err = config.ReadInConfig()
	if err != nil {
		log.Println("error on parsing default configuration file", err)
		return err
	}

	envConfig := viper.New()
	envConfig.SetConfigType("yaml")
	envConfig.AddConfigPath(cfgPath)
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
			Name:        value(config, "app.name", "APP_NAME", "Gin demo"),
			Description: value(config, "app.description", "APP_DESCRIPTION", ""),
			Profile:     value(config, "app.profile", "APP_PROFILE", "dev"),
			CfgPath:     cfgPath,
		},
		DataSource: &DataSource{
			URL:     value(config, "datasource.url", "APP_DATASOURCE_URL", ""),
			Driver:  value(config, "datasource.driver", "APP_DATASOURCE_DRIVER", ""),
			PoolMax: config.GetInt("datasource.pool.max"),
			PoolMin: config.GetInt("datasource.pool.min"),
		},
		Server: &Server{
			Port: config.GetString("server.port"),
		},
		Logging: &Logging{
			Level:  value(config, "logging.level", "APP_LOGGING_LEVEL", "debug"),
			Format: config.GetString("logging.format"),
			Output: config.GetString("logging.output"),
		},
		OTEL: &OTEL{
			ServiceName:  value(config, "otel.service.name", "SERVICE_NAME", "gin-demo"),
			ServiceVer:   value(config, "otel.service.version", "SERVICE_VER", "0.0.0"),
			Insecure:     value(config, "otel.insecure", "INSECURE_MODE", "true"),
			OTLPEndpoint: value(config, "otel.exporter.endpoint", "OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			Logging:      config.GetBool("otel.exporter.logging"),
			Tracer:       config.GetBool("otel.exporter.tracer"),
			Metric:       config.GetBool("otel.exporter.metrics"),
		},
		Security: &Security{
			JWT: &JWT{
				PrivateKey:                      value(config, "security.jwt.private_key", "", ""),
				PublicKey:                       value(config, "security.jwt.public_key", "", ""),
				ExpirationTimeMinutes:           intValue(config, "security.jwt.expiration", "JWT_EXPIRATION", 10),
				RefreshTokenExpirationTimeHours: intValue(config, "security.jwt.refresh_token_expiration", "JWT_REFRESH_TOKEN_EXPIRATION", 72),
				TokenBlackListEnabled:           true,
				TokenBlackListTTL:               intValue(config, "security.jwt.token_blacklist_ttl", "JWT_TOKEN_BLACKLIST_TTL", 30),
			},
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
