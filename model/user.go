package model

import "time"

type User struct {
	Id             string    `xorm:"pk" json:"id" update:"fixed"`
	Email          string    `json:"email" binding:"required,email"`
	PasswordDigest string    `json:"-"`
	Name           string    `json:"name" binding:"required"`
	CreateTime     time.Time `json:"create_time" gorm:"column:created_at"`
	UpdateTime     time.Time `json:"update_time" gorm:"column:updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
