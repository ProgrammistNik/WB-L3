package producer

import (
	"context"

	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/zlog"
)

type KafkaServiceProducer struct {
	producer *kafka.Producer
}

func NewKafkaServiceProducer(producer *kafka.Producer) *KafkaServiceProducer {
	return &KafkaServiceProducer{producer: producer}
}

func (k *KafkaServiceProducer) Send(ctx context.Context, key, value []byte) error {
	if err := k.producer.Send(ctx, key, value); err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to send Kafka message")
		return err
	}
	return nil
}
