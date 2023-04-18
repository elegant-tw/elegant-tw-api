package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	// Debug Mode Enabled
	DebugMode bool `mapstructure:"DEBUG"`

	// Gin Mode
	GinMode string `mapstructure:"GIN_MODE"`

	// Database
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`

	// HTTP Server Port
	ServerAddr string `mapstructure:"SERVER_ADDR"`
	ServerPort string `mapstructure:"SERVER_PORT"`

	// Rate Limit
	RateLimitEnabled     bool   `mapstructure:"RATE_LIMIT_ENABLED"`
	RateLimitFormatted   string `mapstructure:"RATE_LIMIT_FORMATTED"`
	RateLimitStoreRedis  bool   `mapstructure:"RATE_LIMIT_STORE_REDIS"`
	RateLimitRedisHost   string `mapstructure:"RATE_LIMIT_REDIS_HOST"`
	RateLimitRedisPort   string `mapstructure:"RATE_LIMIT_REDIS_PORT"`
	RateLimitRedisPrefix string `mapstructure:"RATE_LIMIT_REDIS_PREFIX"`

	// CORS
	CORSAllowAllOrigin bool   `mapstructure:"CORS_ALLOW_ALL_ORIGINS"`
	CORSAllowOrigins   string `mapstructure:"CORS_ALLOW_ORIGINS"`
}

func Read() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.BindEnv("DEBUG")
	viper.SetDefault("DEBUG", false)

	viper.BindEnv("GIN_MODE")
	viper.SetDefault("GIN_MODE", "release")

	viper.BindEnv("DB_USERNAME")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.SetDefault("DB_PORT", 5432)
	viper.BindEnv("DB_NAME")
	viper.SetDefault("DB_NAME", "elegant")

	viper.BindEnv("SERVER_ADDR")
	viper.SetDefault("SERVER_ADDR", "0.0.0.0")
	viper.BindEnv("SERVER_PORT")
	viper.SetDefault("SERVER_PORT", 3000)

	viper.BindEnv("RATE_LIMIT_ENABLED")
	viper.SetDefault("RATE_LIMIT_ENABLED", false)
	viper.BindEnv("RATE_LIMIT_FORMATTED")
	viper.SetDefault("RATE_LIMIT_FORMATTED", "10-S")
	viper.BindEnv("RATE_LIMIT_STORE_REDIS")
	viper.SetDefault("RATE_LIMIT_STORE_REDIS", false)
	viper.BindEnv("RATE_LIMIT_REDIS_HOST")
	viper.BindEnv("RATE_LIMIT_REDIS_PORT")
	viper.SetDefault("RATE_LIMIT_REDIS_PORT", 6379)
	viper.BindEnv("RATE_LIMIT_REDIS_PREFIX")
	viper.SetDefault("RATE_LIMIT_REDIS_PREFIX", "elegant")

	viper.BindEnv("CORS_ALLOW_ALL_ORIGINS")
	viper.BindEnv("CORS_ALLOW_ORIGINS")

	viper.SetConfigType("dotenv")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("Warning, there is no config file: %v\n", err)
	}

	// Unmarshal the configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
