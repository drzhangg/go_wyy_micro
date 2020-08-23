package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go_wyy_micro/common/conf"
)

var (
	Viper *viper.Viper
)

func main() {
	Viper = conf.InitConf()
	//fmt.Println(v.Get("etcd.port"))
	fmt.Println(Viper.Get("etcd.addr"))
}
