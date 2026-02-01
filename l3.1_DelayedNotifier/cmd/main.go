package main

import (
	"time"

	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/zlog"

	"github.com/ProgrammistNik/WB-L3/l3.1/internal/config"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/handler"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/queue/consumer"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/queue/producer"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/service"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/storage"
)

func main() {
	zlog.Init()
	zlog.Logger.Info().Msg("Starting application...")

	cfgLoader := config.New()
	if err := cfgLoader.Load("config/config.yaml"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to load config")
	}

	var cfg config.Config
	if err := cfgLoader.Unmarshal(&cfg); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to unmarshal config")
	}

	conn, err := rabbitmq.Connect(cfg.RabbitMQ.URL, 3, time.Second)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to create channel")
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"notifications", 
		true,            
		false,           
		false,           
		false,           
		nil,            
	)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to declare queue")
	}

	pub := rabbitmq.NewPublisher(ch, "")
	prod := producer.New(pub, cfg.RabbitMQ.Queue, cfg.RabbitMQ.Retry)

	store := storage.New()
	svc := service.New(prod, store)

	consCfg := rabbitmq.NewConsumerConfig(cfg.RabbitMQ.Queue)
	rmqConsumer := rabbitmq.NewConsumer(ch, consCfg)
	cons := consumer.New(rmqConsumer, svc)
	cons.Start()

	h := handler.New(svc)

	zlog.Logger.Info().Str("addr", cfg.Server.Address).Msg("starting server")
	if err := h.Router().Run(cfg.Server.Address); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("server failed to start")
	}
}
