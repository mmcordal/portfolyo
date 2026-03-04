package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Reminder struct {
	bun.BaseModel `bun:"table:reminders,alias:r"`
	CoreModel

	UserID     int64     `bun:",notnull"`
	User       *User     `bun:"rel:belongs-to,join:user_id=id"`
	Title      string    `bun:",notnull"`
	ReminderAt time.Time `bun:",notnull"`
}
