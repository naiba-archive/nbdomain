package panel

import (
	"github.com/jinzhu/gorm"
)

//Domain 域名
type Domain struct {
	gorm.Model
	UserID  uint
	CatID   uint
	PanelID uint
	Domain  string `gorm:"type:varchar(100);unique_index"`
	Desc    string

	User  User  `json:",omitempty"`
	Panel Panel `json:",omitempty"`
	Cat   Cat   `json:",omitempty"`
}
