package model

import (
	"github.com/jinzhu/gorm"
)

//Order 订单
type Order struct {
	gorm.Model
	What   string
	UserID uint
	Finish bool

	User User `json:",omitempty"`
}
