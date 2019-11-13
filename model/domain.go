package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Domain 域名
type Domain struct {
	gorm.Model
	UserID      uint
	CatID       uint
	PanelID     uint
	Domain      string    `gorm:"type:varchar(100);unique_index"`
	Create      time.Time //注册时间
	Expire      time.Time //到期时间
	WhoisUpdate time.Time
	Cost        int       //购入成本
	Renew       int       //续费成本
	Buy         time.Time //购入时间
	Registrar   string    //注册平台
	Desc        string

	User  User  `json:",omitempty"`
	Panel Panel `json:",omitempty"`
	Cat   Cat   `json:",omitempty"`
}
