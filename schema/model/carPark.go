package model

import "gorm.io/gorm"

type CarPark struct {
	gorm.Model
	Location string
	Number   int32
}

func (u *CarPark) TableName() string {
	return "CarPark"
}
