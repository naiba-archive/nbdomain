package model

//Cat 分类
type Cat struct {
	Common
	PanelID uint64
	UserID  uint64
	Name    string
	NameEn  string
	Index   int

	Domains []Domain `json:",omitempty"`
	Panel   *Panel   `json:",omitempty"`
	User    *User    `json:",omitempty"`
}
