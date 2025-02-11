package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"phonon/pkg/config"
	"phonon/pkg/converter"
	"phonon/pkg/instrumentation"
	"phonon/pkg/queue"
	"phonon/pkg/repository"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.Initialize()
	instrumentation.InitializeLogging()

	logrus.Info("Starting background with configuration:", viper.AllKeys(), viper.AllSettings())

	db, err := repository.NewDatabase()
	if err != nil {
		logrus.Fatal(err)
	}

	audioConverter := converter.NewFFMPEG(viper.GetString("converter.target_format"))

	consumer, err := queue.NewKafkaConsumer(queue.KafkaConfig{
		Brokers: viper.GetStringSlice("mq.kafka.brokers"),
		GroupID: viper.GetString("mq.kafka.audio_conversion.group"),
		Topic:   viper.GetString("mq.kafka.audio_conversion.topic"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	defer consumer.Close()

	audioConversionQueue := queue.NewAudioConversion(audioConverter, db, queue.AudioConversionWithConsumer(consumer))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		audioConversionQueue.StartConsuming(ctx)
	}()

	logrus.Info("Cleanup consumer service started")

	<-stop
	logrus.Info("\nShutting down gracefully...")

	cancel()

	logrus.Info("Cleanup consumer service stopped cleanly.")
}
