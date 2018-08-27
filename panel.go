package panel

import (
	"github.com/spf13/viper"
)

//Config 全局配置
type Config struct {
	Web struct {
		Addr string
	}
	Mail struct {
		SMTP string `mapstructure:"smtp"`
		TLS  bool   `mapstructure:"tls"`
		User string
		Pass string
	}
	Captcha string
}

//CF 全局配置
var CF Config

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.Unmarshal(&CF)
}
