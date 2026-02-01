package consumer

import (
	"context"
	"encoding/json"

	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/service"
	wbfkafka "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type ImageProcessorConsumer struct {
	consumer *wbfkafka.Consumer
	service  *service.Service
	strategy retry.Strategy
}

func NewImageProcessorConsumer(brokers []string, topic, groupID string, service *service.Service, strategy retry.Strategy) *ImageProcessorConsumer {
	return &ImageProcessorConsumer{
		consumer: wbfkafka.NewConsumer(brokers, topic, groupID),
		service:  service,
		strategy: strategy,
	}
}

// StartConsuming запускает цикл обработки сообщений
func (c *ImageProcessorConsumer) StartConsuming(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			zlog.Logger.Info().Msg("Kafka consumer shutting down...")
			return
		default:
			msg, err := c.consumer.FetchWithRetry(ctx, c.strategy)
			if err != nil {
				zlog.Logger.Error().Err(err).Msg("Error fetching Kafka message")
				continue
			}

			var task struct {
				ImageID string `json:"image_id"`
				Path    string `json:"path"`
			}

			if err := json.Unmarshal(msg.Value, &task); err != nil {
				zlog.Logger.Error().Err(err).Msg("Error unmarshalling Kafka message")
				continue
			}

			if err := c.service.ProcessImage(ctx, task.ImageID, task.Path); err != nil {
				zlog.Logger.Error().Err(err).Msg("Error processing image")
				continue
			}

			if err := c.consumer.Commit(ctx, msg); err != nil {
				zlog.Logger.Error().Err(err).Msg("Error committing Kafka message")
			}
		}
	}
}

func (c *ImageProcessorConsumer) Close() error {
	return c.consumer.Close()
}
