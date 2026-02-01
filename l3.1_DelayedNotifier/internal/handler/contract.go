package handler

import (
	"github.com/ProgrammistNik/WB-L3/l3.1/internal/model"
	"github.com/wb-go/wbf/ginext"
)

type Service interface {
	CreateNotification(*ginext.Context, model.Notification) error
	GetStatusByID(*ginext.Context, string) (model.Notification, error)
	DeleteNotify(*ginext.Context, string) error
}
