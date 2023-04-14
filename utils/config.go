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

	// Sentry
	SentryEnabled    bool    `mapstructure:"SENTRY_ENABLED"`
	SentryDSN        string  `mapstructure:"SENTRY_DSN"`
	SentryTracesRate float64 `mapstructure:"SENTRY_TRACESRATE"`
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

	viper.BindEnv("SENTRY_ENABLED")
	viper.SetDefault("SENTRY_ENABLED", false)
	viper.BindEnv("SENTRY_DSN")
	viper.BindEnv("SENTRY_TRACESRATE")
	viper.SetDefault("SENTRY_TRACESRATE", 1.0)

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
