package model

import "gorm.io/gorm"

type CarPark struct {
	gorm.Model
	Location     string `gorm:"uniqueIndex:location_number;size:100"`
	Number       int32  `gorm:"uniqueIndex:location_number"`
	ParkState    int32
	TimeInterval map[int64]TimeInterval
}

type TimeInterval struct {
	StartTime int64
	EndTime   int64
}

func (u *CarPark) TableName() string {
	return "CarPark"
}
