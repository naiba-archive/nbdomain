package panel

import (
	"github.com/jinzhu/gorm"
)

//Domain 域名
type Domain struct {
	gorm.Model
	PanelID uint
	UserID  uint
	CatID   uint
	Normal  string `gorm:"type:varchar(100);unique_index"`
	Custom  string
	Desc    string

	User  User
	Panel Panel
	Cat   Cat
}
