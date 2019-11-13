package panel

import (
	"fmt"
	"time"

	"github.com/naiba/com"

	"github.com/jinzhu/gorm"
)

//User model
type User struct {
	gorm.Model
	Mail           string `gorm:"type:varchar(50);unique_index"`
	Password       string `json:"-"`
	IsAdmin        bool
	Token          *string `gorm:"type:varchar(100);unique_index"`
	UcenterID      string  `gorm:"unique_index"`
	UcenterExtra   string  `gorm:"type:varchar(1000);"`
	Name           string
	Phone          string
	QQ             string
	Weixin         string
	GoldVIPExpire  time.Time
	SuperVIPExpire time.Time

	Domains []Domain `json:",omitempty"`
	Panels  []Panel  `json:",omitempty"`
	Cats    []Cat    `json:",omitempty"`
	Offers  []Offer  `json:",omitempty"`
}

//GenerateToken 生成用户Token
func (u *User) GenerateToken(db *gorm.DB) error {
	token := com.MD5(fmt.Sprintf("%s%d%d%s", u.UcenterID, u.ID, time.Now().UnixNano(), com.RandomString(10)))
	u.Token = &token
	return db.Save(u).Error
}

//SID 字符串ID
func (u User) SID() string {
	return fmt.Sprintf("%d", u.ID)
}
