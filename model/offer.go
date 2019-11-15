package model

//Offer 报价
type Offer struct {
	Common
	Currency string
	Name     string
	Mail     string
	Amount   int
	Domain   string
	UserID   uint

	User User `json:",omitempty"`
}
