package model

//Cat 分类
type Cat struct {
	Common
	PanelID uint
	UserID  uint
	Name    string
	NameEn  string
	Index   int

	Domains []Domain `json:",omitempty"`
	Panel   *Panel   `json:",omitempty"`
	User    *User    `json:",omitempty"`
}
