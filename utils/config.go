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
}

func Read() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

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
