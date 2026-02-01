package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/handler/dto"
	"github.com/ProgrammistNik/WB-L3/l3.5_EventBooker/internal/model"
	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service Service
}

func New(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Router() *ginext.Engine {
	router := ginext.New("release")
	router.Use(ginext.Logger(), ginext.Recovery())

	router.POST("/events", h.CreateEvent)
	router.POST("/events/:id/book", h.BookEvent)
	router.POST("/events/:id/confirm", h.ConfirmBooking)
	router.GET("/events/:id", h.GetEvent)
	router.GET("/events", h.GetEvents)
	router.GET("/events/:id/bookings", h.GetEventBookings)

	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("./web/*.html")

	router.GET("/", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/admin", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "admin.html", nil)
	})
	router.GET("/user", func(c *ginext.Context) {
		c.HTML(http.StatusOK, "user.html", nil)
	})

	return router
}

func (h *Handler) CreateEvent(c *ginext.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid date format, must be RFC3339"})
		return
	}

	event := model.Event{
		Name:       req.Name,
		Date:       date,
		Capacity:   req.Capacity,
		FreeSeats:  req.Capacity,
		PaymentTTL: req.PaymentTTL * 60, // convert minutes → seconds
	}

	if err := h.service.CreateEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.EventResponse{
		ID:        event.ID,
		Name:      event.Name,
		Date:      event.Date.Format(time.RFC3339),
		Capacity:  event.Capacity,
		FreeSeats: event.FreeSeats,
	})
}

func (h *Handler) BookEvent(c *ginext.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid event ID"})
		return
	}

	booking, err := h.service.BookEvent(eventID, req.Seats)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.BookingResponse{
		ID:        booking.ID,
		EventID:   booking.EventID,
		Seats:     booking.Seats,
		Paid:      booking.Paid,
		CreatedAt: booking.CreatedAt.Format(time.RFC3339),
		ExpiresAt: booking.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *Handler) ConfirmBooking(c *ginext.Context) {
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid booking ID"})
		return
	}

	err = h.service.ConfirmBooking(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"status": "confirmed"})
}

func (h *Handler) GetEvent(c *ginext.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid event ID"})
		return
	}

	event, err := h.service.GetEvent(eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, dto.EventResponse{
		ID:        event.ID,
		Name:      event.Name,
		Date:      event.Date.Format(time.RFC3339),
		Capacity:  event.Capacity,
		FreeSeats: event.FreeSeats,
	})
}

func (h *Handler) GetEvents(c *ginext.Context) {
	events, err := h.service.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	resp := make([]dto.EventResponse, 0, len(events))
	for _, e := range events {
		resp = append(resp, dto.EventResponse{
			ID:        e.ID,
			Name:      e.Name,
			Date:      e.Date.Format(time.RFC3339),
			Capacity:  e.Capacity,
			FreeSeats: e.FreeSeats,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetEventBookings(c *ginext.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid event ID"})
		return
	}

	bookings, err := h.service.GetEventBookings(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	resp := make([]dto.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		resp = append(resp, dto.BookingResponse{
			ID:        b.ID,
			EventID:   b.EventID,
			Seats:     b.Seats,
			Paid:      b.Paid,
			CreatedAt: b.CreatedAt.Format(time.RFC3339),
			ExpiresAt: b.ExpiresAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, resp)
}
