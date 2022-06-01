package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone string `gorm:"index;size:100"`
	Name  string
	//身份证
	IdCard string
	//角色
	Role int32
	//审核状态
	CheckState int32
	//账户余额
	Account int64
	//车辆状态
	CarState int32
	//车位Id
	PackId uint
	//车牌号
	CarNumber  string
	UserName   string `gorm:"uniqueIndex;size:100"`
	Password   string
	IsComplete bool
	//预约车位id
	ReserveParkId uint
	//入库时间
	EnterAt int64
}

func (u *User) TableName() string {
	return "User"
}
