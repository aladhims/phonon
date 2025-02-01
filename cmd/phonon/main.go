package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

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

	fileStore := storage.NewLocal("./data")
	audioConverter := converter.NewFFMPEG()

	// Create the service and handler.
	audioService := service.NewAudioService(db, fileStore, audioConverter)

	router := api.NewRouter(audioService)

	// Create the server.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", viper.GetString("server.port")),
		Handler: router,
	}

	// Channel to listen for termination signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	go func() {
		logrus.WithField("addr", server.Addr).Info("starting server")
		if err = server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
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
