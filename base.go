package panel

import (
	"github.com/jinzhu/gorm"
	//MySQL 驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

//Config 全局配置
type Config struct {
	Web struct {
		Addr   string
		Domain string
	}
	Mail struct {
		SMTP string `mapstructure:"smtp"`
		User string
		Pass string
	}
	Database struct {
		User   string
		Pass   string
		Server string
		Name   string
		Loc    string
	}
	ReCaptcha string `mapstructure:"recaptcha"`
}

//CF 全局配置
var CF Config

//DB 数据库连接
var DB *gorm.DB

func init() {
	//加载配置
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.Unmarshal(&CF)
	//连接数据库
	DB, err = gorm.Open("mysql", CF.Database.User+":"+CF.Database.Pass+"@"+CF.Database.Server+"/"+CF.Database.Name+"?charset=utf8&parseTime=True&loc="+CF.Database.Loc)
	if err != nil {
		panic(err)
	}
	//Debug
	DB = DB.Debug()
}
