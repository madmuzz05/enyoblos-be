package config

import (
	"reflect"

	"github.com/spf13/viper"
)

var AppConfig Config

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Port             int    `mapstructure:"APP_PORT"`
	JwtSecret        string `mapstructure:"JWT_SECRET"`
	JwtKey           string `mapstructure:"JWT_KEY"`
	JwtExpiresIn     int64  `mapstructure:"JWT_EXPIRES_IN"`
	DatabaseHost     string `mapstructure:"DB_HOST"`
	DatabasePort     string `mapstructure:"DB_PORT"`
	DatabaseUsername string `mapstructure:"DB_USERNAME"`
	DatabasePassword string `mapstructure:"DB_PASSWORD"`
	DatabaseName     string `mapstructure:"DB_DATABASE"`
	DatabaseSSL      string `mapstructure:"DB_SSL"`
	RateLimitMax     int    `mapstructure:"RATE_LIMIT_MAX"`
	RateLimitWindow  int    `mapstructure:"RATE_LIMIT_WINDOW"`
	RedisHost        string `mapstructure:"REDIS_HOST"`
	RedisPort        string `mapstructure:"REDIS_PORT"`
	RedisPassword    string `mapstructure:"REDIS_PASSWORD"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (err error) {
	// Try to read .env file first
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")        // current directory
	viper.AddConfigPath("./config") // config directory
	viper.AddConfigPath("../")      // parent directory

	// Read config file if exists, ignore if not found
	viper.ReadInConfig()

	// Enable automatic environment variable reading
	// This will override values from .env file if env vars exist
	viper.AutomaticEnv()

	// Auto-bind all struct fields
	bindEnvVars()

	err = viper.Unmarshal(&AppConfig)
	return
}

// bindEnvVars automatically binds all environment variables from struct tags
func bindEnvVars() {
	v := reflect.ValueOf(AppConfig)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("mapstructure"); tag != "" {
			viper.BindEnv(tag)
		}
	}
}
