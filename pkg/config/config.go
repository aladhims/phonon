package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Initialize loads the configuration from config.yaml or overrides from ENV with prefix APP* using Viper.
// Overriding is used to run the service in custom behavior / dependencies i.e. using MySQL instead of SQLite
func Initialize() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file: %v", err)
	}
}
