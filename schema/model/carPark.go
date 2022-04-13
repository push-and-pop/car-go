package model

import "gorm.io/gorm"

type CarPark struct {
	gorm.Model
	Location     string `gorm:"uniqueIndex:location_number;size:100"`
	Number       int32  `gorm:"uniqueIndex:location_number"`
	ParkState    int32
	TimeInterval string
}

type TimeInterval map[int64]struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

func (u *CarPark) TableName() string {
	return "CarPark"
}
