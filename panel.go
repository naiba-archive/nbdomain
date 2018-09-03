package panel

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//Panel 米表
type Panel struct {
	gorm.Model
	UserID   uint
	Domain   string `gorm:"type:varchar(100);unique_index"`
	Name     string
	NameEn   string
	Desc     string
	DescEn   string
	Analysis string

	User    User
	Cats    []Cat
	Domains []Domain
}

func (p Panel) SID() string {
	return fmt.Sprintf("%d", p.ID)
}
