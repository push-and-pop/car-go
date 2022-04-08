package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone string `gorm:"index;size:10"`
	Name  string
}

func (u *User) TableName() string {
	return "User"
}
