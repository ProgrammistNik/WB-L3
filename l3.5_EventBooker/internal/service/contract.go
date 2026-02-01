package service

import (
	"context"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/model"
)

type Storage interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	BookEvent(ctx context.Context, eventID, seats int, paymentTTL time.Duration) (*model.Booking, error)
	ConfirmBooking(ctx context.Context, bookingID int) error
	CancelBooking(ctx context.Context, bookingID int) error
	GetEvent(ctx context.Context, eventID int) (*model.Event, error)
	CancelExpiredBookings(context.Context) error
	GetEvents(ctx context.Context) ([]model.Event, error)
	GetEventBookings(ctx context.Context, eventID int) ([]model.Booking, error)
}
