package instrumentation

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitializeLogging configures Logrus based on configuration settings.
func InitializeLogging() {
	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
