package model

import "time"

type Comment struct {
	ID        int64     `db:"id"`
	ParentID  *int64    `db:"parent_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
