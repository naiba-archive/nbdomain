package nbdomain

import (
	"github.com/jinzhu/gorm"
	//MySQL 驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"

	"github.com/naiba/nbdomain/model"
)

//CF 全局配置
var CF model.Config

//DB 数据库连接
var DB *gorm.DB

func init() {
	//加载配置
	viper.SetConfigName("config")
	viper.AddConfigPath("./data")
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
	if CF.Debug {
		DB = DB.Debug()
	}
	//禁止软删除
	DB = DB.Unscoped()
}
