package model

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	CoreModel

	Name      string       `bun:",notnull"`
	Surname   string       `bun:",notnull"`
	Email     string       `bun:",unique,notnull"`
	Password  string       `bun:",notnull"`
	Assets    []*UserAsset `bun:"rel:has-many,join:id=user_id"`
	Reminders []*Reminder  `bun:"rel:has-many,join:id=user_id"`
}

func (u *User) String() string {
	return u.Name + " " + u.Surname
}
