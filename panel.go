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

	User    User     `json:",omitempty"`
	Cats    []Cat    `json:",omitempty"`
	Domains []Domain `json:",omitempty"`
}

//SID 字符串ID
func (p Panel) SID() string {
	return fmt.Sprintf("%d", p.ID)
}
