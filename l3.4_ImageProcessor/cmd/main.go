package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/config"
	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/handler"
	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/kafka/consumer"
	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/kafka/producer"
	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/service"
	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/storage"
	wbfkafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.Init()
	zlog.Logger.Info().Msg("Starting ImageProcessor service...")

	var cfg config.Config
	cfgLoader := config.New()
	if err := cfgLoader.Load("config/config.yaml"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to load config")
	}
	if err := cfgLoader.Unmarshal(&cfg); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	db := storage.InitDB(cfg.DB)
	defer storage.CloseDB(db)
	store := storage.New(db)

	kafkaProd := wbfkafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	serviceProducer := producer.NewKafkaServiceProducer(kafkaProd)

	storagePath := cfg.StoragePath

	srv := service.New(store, serviceProducer, storagePath)

	h := handler.New(srv)
	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: h.Router(),
	}

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  1,
	}
	kafkaConsumer := consumer.NewImageProcessorConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, srv, strategy)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск Kafka consumer
	go kafkaConsumer.StartConsuming(ctx)
	defer kafkaConsumer.Close()

	// Запуск HTTP сервера
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Logger.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	zlog.Logger.Info().Msg("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("HTTP server forced to shutdown")
	}

	cancel()
	time.Sleep(time.Second)
}
