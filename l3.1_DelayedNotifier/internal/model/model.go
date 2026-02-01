package model

import "time"

type Notification struct {
	ID        string    `json:"id"` // id самого уведомления
	Message   string    `json:"message"`
	SendAt    time.Time `json:"send_at"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}
