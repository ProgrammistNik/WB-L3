package model

import "time"

type URL struct {
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreateAt    time.Time `json:"create_at"`
}

type Click struct {
	ID        int       `json:"id"`
	LinkID    int       `json:"link_id"`
	UserAgent string    `json:"user_agent"`
	CreateAt  time.Time `json:"create_at"`
}

type AnalyticsResult struct {
	Group string `json:"group"`
	Count int    `json:"count"`
}
