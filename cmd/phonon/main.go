package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"phonon/pkg/queue"
	"syscall"

	"phonon/pkg/api"
	"phonon/pkg/config"
	"phonon/pkg/converter"
	"phonon/pkg/instrumentation"
	"phonon/pkg/repository"
	"phonon/pkg/service"
	"phonon/pkg/storage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.Initialize()
	instrumentation.InitializeLogging()

	// Initialize dependencies.
	db, err := repository.NewDatabase()
	if err != nil {
		logrus.Fatal(err)
	}

	filestore, err := storage.NewFilestore(storage.Config{
		Type:     storage.StorageType(viper.GetString("storage.type")),
		BasePath: viper.GetString("storage.local.base_path")})
	if err != nil {
		logrus.Fatal(err)
	}

	audioConverter := converter.NewFFMPEG()

	// Initialize Kafka producer
	producer, err := queue.NewKafkaProducer(queue.KafkaConfig{
		Brokers: viper.GetStringSlice("mq.kafka.brokers"),
		Topic:   viper.GetString("mq.kafka.topic"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer producer.Close()

	// Create the service and handler.
	audioService := service.NewAudioService(db, filestore, audioConverter, producer)

	router := api.NewRouter(audioService, producer)

	// Create the server.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", viper.GetString("server.port")),
		Handler: router,
	}

	// Channel to listen for termination signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logrus.WithField("addr", server.Addr).Info("starting server")
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for a termination signal.
	<-stop
	logrus.Info("\nShutting down gracefully...")

	// Create a context with a timeout to allow ongoing requests to complete.
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("server.shutdown_timeout"))
	defer cancel()

	// Gracefully shutdown the server.
	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server shutdown failed: %v", err)
	}

	logrus.Info("Server stopped cleanly.")
}
