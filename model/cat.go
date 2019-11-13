package model

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

	Domains []Domain `json:",omitempty"`
	Panel   *Panel   `json:",omitempty"`
	User    *User    `json:",omitempty"`
}
