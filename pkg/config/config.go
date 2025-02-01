package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Initialize loads the configuration from config.yaml using Viper.
func Initialize() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file: %v", err)
	}
}
