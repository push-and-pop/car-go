package util

import (
	"car-go/schema/model"

	"github.com/spf13/viper"
)

func ReadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../config/.")
	viper.AddConfigPath("./config/.")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func HasNoIntersection(interval model.TimeInterval, start, end int64) bool {
	for _, value := range interval {
		if start >= value.EndTime || end <= value.StartTime {
			continue
		}
		return false
	}
	return true
}
