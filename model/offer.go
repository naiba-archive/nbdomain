package model

//Offer 报价
type Offer struct {
	Common
	Domain   string `json:"domain,omitempty"`
	Currency string `json:"currency,omitempty"`
	Amount   int    `json:"amount,omitempty"`
	Name     string `json:"name,omitempty"`
	Mail     string `json:"mail,omitempty"`

	UserID uint64 `json:"user_id,omitempty"`
	User   User   `json:"-"`
}
