package model

import (
	"database/sql"
	"time"
)

type ImageStatus string

const (
	StatusPending    ImageStatus = "pending"
	StatusProcessing ImageStatus = "processing"
	StatusCompleted  ImageStatus = "completed"
	StatusFailed     ImageStatus = "failed"
)

type Image struct {
	ID            string         `json:"id"`
	OriginalPath  string         `json:"original_path"`
	ResizedPath   sql.NullString `json:"resized_path,omitempty"`
	ThumbPath     sql.NullString `json:"thumb_path,omitempty"`
	WatermarkPath sql.NullString `json:"watermark_path,omitempty"`
	Status        ImageStatus    `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	ProcessedAt   sql.NullTime   `json:"processed_at"`
}

func (i *Image) GetResizedPath() string {
	if i.ResizedPath.Valid {
		return i.ResizedPath.String
	}
	return ""
}

func (i *Image) GetThumbPath() string {
	if i.ThumbPath.Valid {
		return i.ThumbPath.String
	}
	return ""
}

func (i *Image) GetWatermarkPath() string {
	if i.WatermarkPath.Valid {
		return i.WatermarkPath.String
	}
	return ""
}

func (i *Image) GetProcessedAt() *time.Time {
	if i.ProcessedAt.Valid {
		return &i.ProcessedAt.Time
	}
	return nil
}
