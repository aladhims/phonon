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

	logrus.Info("Starting phonon with configuration:", viper.AllKeys(), viper.AllSettings())

	db, err := repository.NewDatabase()
	if err != nil {
		logrus.Fatal(err)
	}

	filestore, err := storage.NewFilestore(storage.Config{
		Type:     storage.Type(viper.GetString("storage.type")),
		BasePath: viper.GetString("storage.local.base_path")})
	if err != nil {
		logrus.Fatal(err)
	}

	audioConverter := converter.NewFFMPEG(viper.GetString("converter.target_format"))

	producer, err := queue.NewKafkaProducer(queue.KafkaConfig{
		Brokers: viper.GetStringSlice("mq.kafka.brokers"),
		Topic:   viper.GetString("mq.kafka.audio_conversion.topic"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer producer.Close()

	audioConversionQueue := queue.NewAudioConversion(audioConverter, db, queue.AudioConversionWithProducer(producer))

	audioService := service.NewAudioService(db, filestore, audioConversionQueue)

	router := api.NewRouter(audioService, producer)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", viper.GetString("server.port")),
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logrus.WithField("addr", server.Addr).Info("starting server")
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Server failed: %v", err)
		}
	}()

	<-stop
	logrus.Info("\nShutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("server.shutdown_timeout"))
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server shutdown failed: %v", err)
	}

	logrus.Info("Server stopped cleanly.")
}
