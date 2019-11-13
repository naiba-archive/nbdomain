package model

import (
	"github.com/jinzhu/gorm"
)

//Offer 报价
type Offer struct {
	gorm.Model
	Currency string
	Name     string
	Mail     string
	Amount   int
	Domain   string
	UserID   uint

	User User `json:",omitempty"`
}
