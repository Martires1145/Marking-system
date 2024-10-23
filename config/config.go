package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("err:%#v", err)
		return
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(dir + "/config")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("err:%#v", err)
		return
	}
}
