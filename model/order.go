package model

//Order 订单
type Order struct {
	Common
	What   string
	UserID uint64
	Finish bool

	User User `json:",omitempty"`
}
