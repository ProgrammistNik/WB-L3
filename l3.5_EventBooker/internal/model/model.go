package model

import "time"

type Event struct {
	ID         int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Date       time.Time `db:"date" json:"date"`
	Capacity   int       `db:"capacity" json:"capacity"`
	FreeSeats  int       `db:"free_seats" json:"freeSeats"`
	PaymentTTL int       `db:"payment_ttl" json:"paymentTTL"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
}

type Booking struct {
	ID        int       `json:"id"`
	EventID   int       `json:"eventId"`
	Seats     int       `json:"seats"`
	Paid      bool      `json:"paid"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
