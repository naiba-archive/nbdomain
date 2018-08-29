package panel

import (
	"github.com/jinzhu/gorm"
)

//Panel 米表
type Panel struct {
	gorm.Model
	UserID         uint
	Domain         string `gorm:"type:varchar(100);unique_index"`
	Name           string
	NameEn         string
	Logo           string
	LogoEn         string
	Desc           string
	DescEn         string
	GoogleAnalysis string

	User    User
	Cats    []Cat
	Domains []Domain
}
