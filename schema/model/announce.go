package model

import "gorm.io/gorm"

type Announce struct {
	gorm.Model
	Msg string
}

func (a *Announce) TableName() string {
	return "Announce"
}
