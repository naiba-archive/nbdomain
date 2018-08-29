package panel

import (
	"github.com/jinzhu/gorm"
)

//Cat 分类
type Cat struct {
	gorm.Model
	PanelID uint
	UserID  uint
	Name    string
	NameEn  string
	Index   int

	Domains []Domain
	Panel   Panel
	User    User
}
