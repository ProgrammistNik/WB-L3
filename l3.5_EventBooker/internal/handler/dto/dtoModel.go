package dto

type CreateEventRequest struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	Capacity   int    `json:"capacity"`
	PaymentTTL int    `json:"paymentTTL"`
}

type EventResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Date      string `json:"date"`
	Capacity  int    `json:"capacity"`
	FreeSeats int    `json:"freeSeats"`
}

type CreateBookingRequest struct {
	Seats int `json:"seats"`
}

type BookingResponse struct {
	ID        int    `json:"id"`
	EventID   int    `json:"eventId"`
	Seats     int    `json:"seats"`
	Paid      bool   `json:"paid"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
}
