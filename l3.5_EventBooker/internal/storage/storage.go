package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/model"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type Storage struct{ db *dbpg.DB }

func New(db *dbpg.DB) *Storage { return &Storage{db: db} }

func (st *Storage) CreateEvent(ctx context.Context, e *model.Event) error {
	tx, err := st.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO events (name, date, capacity, free_seats, payment_ttl, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id
	`, e.Name, e.Date, e.Capacity, e.FreeSeats, e.PaymentTTL, e.CreatedAt, e.UpdatedAt)

	if err = row.Scan(&e.ID); err != nil {
		return err
	}

	return tx.Commit()
}

func (st *Storage) GetEvents(ctx context.Context) ([]model.Event, error) {
	rows, err := st.db.Master.QueryContext(ctx, `
		SELECT id, name, date, capacity, free_seats, payment_ttl, created_at, updated_at
		FROM events ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event

	for rows.Next() {
		var e model.Event
		err := rows.Scan(
			&e.ID, &e.Name, &e.Date, &e.Capacity,
			&e.FreeSeats, &e.PaymentTTL, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (st *Storage) GetEvent(ctx context.Context, id int) (*model.Event, error) {
	var e model.Event
	err := st.db.Master.QueryRowContext(ctx, `
		SELECT id, name, date, capacity, free_seats, payment_ttl, created_at, updated_at
		FROM events WHERE id=$1
	`, id).Scan(
		&e.ID, &e.Name, &e.Date, &e.Capacity, &e.FreeSeats,
		&e.PaymentTTL, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (st *Storage) BookEvent(ctx context.Context, eventId, seats int, ttl time.Duration) (*model.Booking, error) {
	tx, err := st.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var free int
	err = tx.QueryRowContext(ctx, `SELECT free_seats FROM events WHERE id=$1 FOR UPDATE`, eventId).
		Scan(&free)

	if err != nil {
		return nil, fmt.Errorf("event not found")
	}

	if free < seats {
		return nil, fmt.Errorf("not enough free seats")
	}

	now := time.Now()
	b := &model.Booking{
		EventID:   eventId,
		Seats:     seats,
		Paid:      false,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}

	err = tx.QueryRowContext(ctx, `
		INSERT INTO bookings (event_id, seats, paid, created_at, expires_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`, b.EventID, b.Seats, b.Paid, b.CreatedAt, b.ExpiresAt).
		Scan(&b.ID)

	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE events SET free_seats = free_seats - $1 WHERE id=$2
	`, seats, eventId)
	if err != nil {
		return nil, err
	}

	return b, tx.Commit()
}

func (st *Storage) ConfirmBooking(ctx context.Context, id int) error {
	_, err := st.db.Master.ExecContext(ctx,
		`UPDATE bookings SET paid=true WHERE id=$1`, id)
	return err
}

func (st *Storage) CancelExpiredBookings(ctx context.Context) error {
	rows, err := st.db.Master.QueryContext(ctx, `
		SELECT id FROM bookings WHERE paid=false AND expires_at <= now()
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for _, id := range ids {
		zlog.Logger.Warn().Int("booking", id).Msg("auto-cancel expired booking")

		if err := st.CancelBooking(ctx, id); err != nil {
			zlog.Logger.Error().
				Err(err).
				Int("booking", id).
				Msg("failed to cancel expired booking")
		}
	}

	return nil
}

func (st *Storage) CancelBooking(ctx context.Context, id int) error {
	tx, err := st.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var eventId, seats int
	err = tx.QueryRowContext(ctx, `
		SELECT event_id, seats FROM bookings WHERE id=$1
	`, id).Scan(&eventId, &seats)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM bookings WHERE id=$1
	`, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE events SET free_seats = free_seats + $1 WHERE id=$2
	`, seats, eventId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (st *Storage) GetEventBookings(ctx context.Context, eventId int) ([]model.Booking, error) {
	rows, err := st.db.Master.QueryContext(ctx, `
		SELECT id, event_id, seats, paid, created_at, expires_at
		FROM bookings WHERE event_id=$1 ORDER BY id
	`, eventId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Booking
	for rows.Next() {
		var b model.Booking
		err := rows.Scan(&b.ID, &b.EventID, &b.Seats, &b.Paid, &b.CreatedAt, &b.ExpiresAt)
		if err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}
