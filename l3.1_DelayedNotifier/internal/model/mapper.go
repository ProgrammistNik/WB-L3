package model

import (
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.1/internal/handler/dto"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/zlog"
)


func CastToNotification(request dto.NotificationRequest) Notification {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to load Moscow timezone, falling back to UTC")
		location = time.UTC
	}

	mskTime := request.SendAt.In(location)

	zlog.Logger.Info().
		Str("send_at_utc", request.SendAt.Format("2006-01-02 15:04:05 UTC")).
		Str("send_at_msk", mskTime.Format("2006-01-02 15:04:05")).
		Msg("Notification scheduled")

	return Notification{
		ID:        uuid.New().String(),
		Message:   request.Message,
		SendAt:    mskTime, 
		CreatedAt: time.Now().UTC(),
		Status:    "scheduled",
	}
}
