package service

import (
	"fmt"

	"github.com/ProgrammistNik/WB-L3/l3.1/internal/model"
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/queue/producer"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type Service struct {
	producer *producer.Producer
	storage  Storage
}

func New(p *producer.Producer, st Storage) *Service {
	return &Service{
		producer: p,
		storage:  st,
	}
}

func (s *Service) CreateNotification(ctx *ginext.Context, notification model.Notification) error {
	s.storage.Set(notification)
	err := s.producer.Publish(notification)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("notification_id", notification.ID).
			Msg("failed to publish notification to queue")
		return err
	}
	return nil
}

func (s *Service) ProcessNotification(notification model.Notification) {
	storedNotif, ok := s.storage.Get(notification.ID)
	if !ok {
		zlog.Logger.Warn().
			Str("id", notification.ID).
			Msg("notification not found in storage — skipping")
		return
	}

	if storedNotif.Status == "processed" {
		zlog.Logger.Info().
			Str("id", notification.ID).
			Msg("notification already processed — skipping")
		return
	}

	if storedNotif.Status == "canceled" {
		zlog.Logger.Info().
			Str("id", notification.ID).
			Msg("notification was canceled — skipping processing")
		return
	}

	notification.Status = "processed"
	s.storage.Set(notification)

	zlog.Logger.Info().
		Str("id", notification.ID).
		Msg("notification processed successfully")
}

func (s *Service) GetStatusByID(ctx *ginext.Context, id string) (model.Notification, error) {
	res, ok := s.storage.Get(id)
	if !ok {
		return model.Notification{}, fmt.Errorf("id not found")
	}

	return res, nil
}

func (s *Service) DeleteNotify(ctx *ginext.Context, id string) error {
	notif, ok := s.storage.Get(id)
	if !ok {
		return fmt.Errorf("notification not found")
	}

	notif.Status = "canceled"
	s.storage.Set(notif)
	return nil
}
