package service

import "github.com/ProgrammistNik/WB-L3/l3.1/internal/model"

type Storage interface {
	Set(notif model.Notification)
	Get(id string) (model.Notification, bool)
}
