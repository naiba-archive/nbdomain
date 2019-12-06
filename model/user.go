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
	Avatar      string     `gorm:"-" json:"avatar"`
	Mail        string     `gorm:"type:varchar(50);unique_index" json:"mail,omitempty"`
	Password    string     `json:"-"`
	IsAdmin     bool       `json:"is_admin,omitempty"`
	Token       *string    `gorm:"type:varchar(100);unique_index" json:"token,omitempty"`
	TokenExpire *time.Time `json:"token_expire,omitempty"`
	Name        string     `json:"name,omitempty"`
	Phone       string     `json:"phone,omitempty"`
	QQ          string     `json:"qq,omitempty"`
	Weixin      string     `json:"weixin,omitempty"`

	Domains []Domain `json:"-"`
	Panels  []Panel  `json:"-"`
	Cats    []Cat    `json:"-"`
	Offers  []Offer  `json:"-"`
}

// AfterFind ..
func (u *User) AfterFind() {
	u.Avatar = "https://i.pravatar.cc/150?u=" + com.MD5(u.Mail)
	if u.Name == "" {
		u.Name = "管理员"
	}
}

//GenerateToken 生成用户Token
func (u *User) GenerateToken(db *gorm.DB, autoLogin bool) error {
	token := com.MD5(fmt.Sprintf("%d%d%s", u.ID, time.Now().UnixNano(), com.RandomString(10)))
	u.Token = &token
	expire := time.Now().AddDate(1, 0, 0)
	if !autoLogin {
		expire = time.Now().Add(time.Hour * 2)
	}
	u.TokenExpire = &expire
	return db.Save(u).Error
}

//SID 字符串ID
func (u User) SID() string {
	return fmt.Sprintf("%d", u.ID)
}
