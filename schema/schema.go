package schema

import (
	"car-go/schema/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func LinkDb() {
	var err error
	dsn := viper.GetString("mysql.user") +
		":" +
		viper.GetString("mysql.password") +
		"@tcp(" +
		viper.GetString("mysql.addr") +
		")/" +
		viper.GetString("mysql.db") +
		"?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Db.AutoMigrate(model.User{})
}
