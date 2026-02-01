package model

import "time"

type Item struct {
	ID        int64
	Type      string
	Category  string
	Amount    float64
	Date      time.Time
	CreatedAt time.Time
}

type CreateItemRequest struct {
	Type     string  `json:"type" binding:"required"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount" binding:"required,gte=0"`
	Date     string  `json:"date" binding:"required"`
}

type AnalyticsResponse struct {
	Sum    float64 `json:"sum"`
	Avg    float64 `json:"avg"`
	Count  int64   `json:"count"`
	Median float64 `json:"median"`
	P90    float64 `json:"p90"`
}

type ItemsFilter struct {
	From     *time.Time
	To       *time.Time
	Category *string
	Limit    *int
	Offset   *int
}
