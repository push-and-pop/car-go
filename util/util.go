package util

import "github.com/spf13/viper"

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
