package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"phonon/pkg/config"
	"phonon/pkg/instrumentation"
	"phonon/pkg/queue"
	"phonon/pkg/storage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.Initialize()
	instrumentation.InitializeLogging()

	filestore, err := storage.NewFilestore(storage.Config{
		Type:     storage.StorageType(viper.GetString("storage.type")),
		BasePath: viper.GetString("storage.local.base_path")})
	if err != nil {
		logrus.Fatal(err)
	}

	consumer, err := queue.NewKafkaConsumer(queue.KafkaConfig{
		Brokers: viper.GetStringSlice("mq.kafka.brokers"),
		GroupID: viper.GetString("mq.kafka.group"),
		Topic:   viper.GetString("mq.kafka.topic"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer consumer.Close()

	cleanupHandler := queue.NewCleanupHandler(filestore)

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming janitor messages
	go func() {
		if err := consumer.Consume(ctx, cleanupHandler); err != nil {
			logrus.WithError(err).Error("janitor consumer failed")
		}
	}()

	logrus.Info("Cleanup consumer service started")

	// Wait for termination signal
	<-stop
	logrus.Info("\nShutting down gracefully...")

	// Cancel the context to stop the consumer
	cancel()

	logrus.Info("Cleanup consumer service stopped cleanly.")
}
