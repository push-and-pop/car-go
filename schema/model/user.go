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
}

func (u *User) TableName() string {
	return "User"
}
