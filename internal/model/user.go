package model

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"users,alias:u"`
	CoreModel
	Name     string `json:"name" `
	Surname  string `json:"surname"`
	Email    string `json:"email" bun:",unique"`
	Password string `json:"password"`
}

func (u *User) String() string {
	return u.Name + " " + u.Surname
}
