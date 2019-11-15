package model

import (
	"fmt"
	"time"

	"github.com/naiba/com"

	"github.com/jinzhu/gorm"
)

//User model
type User struct {
	gorm.Model
	Mail        string `gorm:"type:varchar(50);unique_index"`
	Password    string `json:"-"`
	IsAdmin     bool
	Token       *string `gorm:"type:varchar(100);unique_index"`
	TokenExpire time.Time
	Name        string
	Phone       string
	QQ          string
	Weixin      string

	Domains []Domain `json:",omitempty"`
	Panels  []Panel  `json:",omitempty"`
	Cats    []Cat    `json:",omitempty"`
	Offers  []Offer  `json:",omitempty"`
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
