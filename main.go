package main

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	Viper *viper.Viper
)

func main() {
	//Viper = conf.InitConf()
	////fmt.Println(v.Get("etcd.port"))
	//fmt.Println(Viper.Get("etcd.addr"))

	var arr []string = []string{"1","2"}
	fmt.Println(arr)


	var ma = make(map[string][]string)

	ma["one"] = []string{"1","2"}

	mm := map[string][]string{
		"one":[]string{"12","23"},
	}

	fmt.Println(ma)
	fmt.Println(mm)
}
