package producer

import (
	"encoding/json"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.1/internal/config"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/model"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type Producer struct {
	publisher *rabbitmq.Publisher
	queueName string
	retryCfg  config.RabbitRetryConfig
}

func New(p *rabbitmq.Publisher, queueName string, retryCfg config.RabbitRetryConfig) *Producer {
	return &Producer{
		publisher: p,
		queueName: queueName,
		retryCfg:  retryCfg,
	}
}

func (p *Producer) Publish(n model.Notification) error {
	body, err := json.Marshal(&n)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("notification_id", n.ID).
			Msg("failed to marshal notification")
		return err
	}

	return p.publisher.PublishWithRetry(
		body,
		p.queueName,
		"application/json",
		retry.Strategy{
			Attempts: p.retryCfg.Attempts,
			Delay:    time.Duration(p.retryCfg.DelayMS) * time.Millisecond,
			Backoff:  p.retryCfg.Backoff,
		},
	)
}
