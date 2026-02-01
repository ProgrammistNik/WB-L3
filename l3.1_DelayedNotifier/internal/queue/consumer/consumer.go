package consumer

import (
	"encoding/json"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.1/internal/model"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/zlog"
)

type Consumer struct {
	consumer *rabbitmq.Consumer
	service  ConsumerService
}

func New(c *rabbitmq.Consumer, s ConsumerService) *Consumer {
	return &Consumer{
		consumer: c,
		service:  s,
	}
}

func (c *Consumer) Start() {
	msgChan := make(chan []byte)

	go func() {
		for msg := range msgChan {
			var n model.Notification
			if err := json.Unmarshal(msg, &n); err != nil {
				zlog.Logger.Error().Err(err).Msg("failed to decode message")
				continue
			}

			delay := time.Until(n.SendAt)
			if delay <= 0 {
				zlog.Logger.Info().
					Str("id", n.ID).
					Msg("sending overdue notification immediately")
				c.service.ProcessNotification(n)
			} else {
				zlog.Logger.Info().
					Str("id", n.ID).
					Dur("delay", delay).
					Msg("delaying processing of notification")
				time.Sleep(delay)
				c.service.ProcessNotification(n)
			}
		}
	}()

	go func() {
		if err := c.consumer.Consume(msgChan); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to consume")
		}
	}()
}
