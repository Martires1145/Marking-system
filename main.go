package main

import (
	"github.com/spf13/viper"
	"marking/common"
	_ "marking/config"
	"marking/route"
)

func main() {
	common.InitMySQL()
	r := route.GetRoute()
	port := viper.GetString("web.port")
	_ = r.Run(port)
}
