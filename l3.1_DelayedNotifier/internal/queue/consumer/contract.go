package consumer

import "github.com/ProgrammistNik/WB-L3/l3.1/internal/model"

type ConsumerService interface {
	ProcessNotification(model.Notification)
}
