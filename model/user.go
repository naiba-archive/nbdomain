package model

import (
	"fmt"
	"time"

	"github.com/naiba/com"

	"github.com/jinzhu/gorm"
)

//User model
type User struct {
	Common
	Mail        string    `gorm:"type:varchar(50);unique_index" json:"mail,omitempty"`
	Password    string    `json:"-"`
	IsAdmin     bool      `json:"is_admin,omitempty"`
	Token       *string   `gorm:"type:varchar(100);unique_index" json:"token,omitempty"`
	TokenExpire time.Time `json:"token_expire,omitempty"`
	Name        string    `json:"name,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	QQ          string    `json:"qq,omitempty"`
	Weixin      string    `json:"weixin,omitempty"`

	Domains []Domain `json:"domains,omitempty"`
	Panels  []Panel  `json:"panels,omitempty"`
	Cats    []Cat    `json:"cats,omitempty"`
	Offers  []Offer  `json:"offers,omitempty"`
}

//GenerateToken 生成用户Token
func (u *User) GenerateToken(db *gorm.DB, autoLogin bool) error {
	token := com.MD5(fmt.Sprintf("%d%d%s", u.ID, time.Now().UnixNano(), com.RandomString(10)))
	u.Token = &token
	if autoLogin {
		u.TokenExpire = time.Now().AddDate(1, 0, 0)
	} else {
		u.TokenExpire = time.Now().Add(time.Hour * 2)
	}
	return db.Save(u).Error
}

//SID 字符串ID
func (u User) SID() string {
	return fmt.Sprintf("%d", u.ID)
}
