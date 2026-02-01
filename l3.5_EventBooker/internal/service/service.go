package service

import (
	"context"
	"errors"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/model"
)

type Service struct {
	storage Storage
}

func New(st Storage) *Service {
	return &Service{storage: st}
}

func (s *Service) CreateEvent(ctx context.Context, event *model.Event) error {
	event.CreatedAt = time.Now()
	event.UpdatedAt = event.CreatedAt
	event.FreeSeats = event.Capacity

	return s.storage.CreateEvent(ctx, event)
}

func (s *Service) BookEvent(eventId, seats int) (*model.Booking, error) {
	event, err := s.storage.GetEvent(context.Background(), eventId)
	if err != nil {
		return nil, errors.New("event not found")
	}

	ttl := time.Duration(event.PaymentTTL) * time.Second

	return s.storage.BookEvent(
		context.Background(),
		eventId,
		seats,
		ttl,
	)
}

func (s *Service) ConfirmBooking(bookingId int) error {
	return s.storage.ConfirmBooking(context.Background(), bookingId)
}

func (s *Service) GetEvent(id int) (*model.Event, error) {
	return s.storage.GetEvent(context.Background(), id)
}

func (s *Service) GetEvents() ([]model.Event, error) {
	return s.storage.GetEvents(context.Background())
}

func (s *Service) GetEventBookings(eventId int) ([]model.Booking, error) {
	return s.storage.GetEventBookings(context.Background(), eventId)
}

func (s *Service) CancelExpiredBookings() error {
	return s.storage.CancelExpiredBookings(context.Background())
}
