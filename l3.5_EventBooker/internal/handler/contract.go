package handler

import (
	"context"

	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/model"
)

type Service interface {
	CreateEvent(context.Context, *model.Event) error
	BookEvent(int, int) (*model.Booking, error)
	ConfirmBooking(int) error
	GetEvent(int) (*model.Event, error)
	GetEvents() ([]model.Event, error)
	GetEventBookings(eventID int) ([]model.Booking, error)
}
