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
	Mail     string `gorm:"type:varchar(50);unique_index"`
	Password string `json:"-"`
	IsAdmin  bool
	Token    string `gorm:"type:varchar(100);unique_index"`
	Expire   time.Time
}

//GenerateToken 生成用户Token
func (u *User) GenerateToken() error {
	u.Token = com.MD5(fmt.Sprintf("%d%d%s", u.ID, time.Now().UnixNano(), com.RandomString(10)))
	return DB.Save(u).Error
}
