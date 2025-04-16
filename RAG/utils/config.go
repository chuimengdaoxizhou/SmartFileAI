package utils

import (
	"github.com/spf13/viper"
)

var Config *viper.Viper

func InitConfig() {
	Config = viper.New()
	Config.SetConfigFile("config/config.yaml")
	if err := Config.ReadInConfig(); err != nil {
		panic("Failed to read config: " + err.Error())
	}
}
