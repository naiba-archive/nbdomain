package model

import (
	"time"
)

//Domain 域名
type Domain struct {
	Common
	UserID      uint64    `json:"user_id,omitempty"`
	CatID       uint64    `json:"cat_id,omitempty"`
	PanelID     uint64    `json:"panel_id,omitempty"`
	Domain      string    `gorm:"type:varchar(100);unique_index" json:"domain,omitempty"`
	Create      time.Time `json:"create,omitempty"` //注册时间
	Expire      time.Time `json:"expire,omitempty"` //到期时间
	WhoisUpdate time.Time `json:"whois_update,omitempty"`
	Cost        int       `json:"cost,omitempty"`      //购入成本
	Renew       int       `json:"renew,omitempty"`     //续费成本
	Buy         time.Time `json:"buy,omitempty"`       //购入时间
	Registrar   string    `json:"registrar,omitempty"` //注册平台
	Desc        string    `json:"desc,omitempty"`

	User  User  `json:"-"`
	Panel Panel `json:"-"`
	Cat   Cat   `json:"-"`
}
