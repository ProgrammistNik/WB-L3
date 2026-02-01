package dto

import "time"

type NotificationRequest struct {
	Message string    `json:"message"`
	SendAt  time.Time `json:"send_at"`
}
