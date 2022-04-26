package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	//订单类型
	Type int32
	//订单状态
	State int32
	//车位id
	PackId uint
	//结束时间
	EndAt int64
}

func (u *Order) TableName() string {
	return "Order"
}
