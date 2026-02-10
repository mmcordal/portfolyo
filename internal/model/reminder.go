package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Reminder struct {
	bun.BaseModel `bun:"reminders,alias:r"`
	CoreModel
	UserID int64     `json:"user_id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
}
