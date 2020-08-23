package conf

import "github.com/spf13/viper"

func init() {
	InitConf()
}

func InitConf() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("./config/config.yml")


	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return v
}

