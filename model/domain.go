package model

import (
	"time"
)

var timeLayouts = []string{
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05-0700",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02",
	"02.01.2006 15:04:05",
	time.RFC1123,     //= "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,    //= "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	time.RFC850,      //= "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC3339,     //= "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, //= "2006-01-02T15:04:05.999999999Z07:00"
}

// ParseWhoisTime ..
func ParseWhoisTime(t string) time.Time {
	for _, layout := range timeLayouts {
		tt, e := time.Parse(layout, t)
		if e == nil {
			return tt
		}
	}
	return time.Now()
}

//Domain 域名
type Domain struct {
	Common
	UserID      uint64     `json:"user_id,omitempty"`
	CatID       uint64     `json:"cat_id,omitempty"`
	PanelID     uint64     `json:"panel_id,omitempty"`
	Domain      string     `gorm:"type:varchar(100);unique_index" json:"domain,omitempty"`
	Create      *time.Time `json:"create,omitempty"` //注册时间
	Expire      *time.Time `json:"expire,omitempty"` //到期时间
	WhoisUpdate *time.Time `json:"whois_update,omitempty"`
	Cost        int        `json:"cost,omitempty"`      //购入成本
	Renew       int        `json:"renew,omitempty"`     //续费成本
	Buy         *time.Time `json:"buy,omitempty"`       //购入时间
	Registrar   string     `json:"registrar,omitempty"` //注册平台
	Desc        string     `json:"desc,omitempty"`

	User  User  `json:"-"`
	Panel Panel `json:"-"`
	Cat   Cat   `json:"-"`
}
