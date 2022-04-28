package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	//订单类型
	Type int32
	//订单状态
	State int32
	//用户id
	UserId uint
	//车位id
	PackId uint
	//开始时间
	StartAt int64
	//结束时间
	EndAt int64
	//价格
	Price int64
}

func (u *Order) TableName() string {
	return "Order"
}
