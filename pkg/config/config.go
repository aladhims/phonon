package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Initialize loads the configuration from config.yaml or overrides from ENV with prefix APP* using Viper.
// Overriding is used to run the service in custom behavior / dependencies i.e. using MySQL instead of SQLite
func Initialize() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	bindEnvVariables()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file: %v", err)
	}
}

// bindEnvVariables binds all configuration keys to their corresponding environment variables
func bindEnvVariables() {
	viper.BindEnv("log.level")

	viper.BindEnv("server.port")
	viper.BindEnv("server.shutdown_timeout")
	viper.BindEnv("server.max_upload_size")

	viper.BindEnv("database.driver")
	viper.BindEnv("database.sqlite.path")
	viper.BindEnv("database.sqlite.seed")
	viper.BindEnv("database.mysql.host")
	viper.BindEnv("database.mysql.port")
	viper.BindEnv("database.mysql.database")
	viper.BindEnv("database.mysql.username")
	viper.BindEnv("database.mysql.password")

	viper.BindEnv("storage.type")
	viper.BindEnv("storage.local.base_path")

	viper.BindEnv("mq.kafka.brokers")
	viper.BindEnv("mq.kafka.audio_conversion.group")
	viper.BindEnv("mq.kafka.audio_conversion.topic")
}
