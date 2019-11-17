package model

//Cat 分类
type Cat struct {
	Common
	PanelID uint64 `json:"panel_id,omitempty"`
	UserID  uint64 `json:"user_id,omitempty"`
	Name    string `json:"name,omitempty"`
	NameEn  string `json:"name_en,omitempty"`
	Index   int    `json:"index,omitempty"`

	Domains []Domain `json:"-"`
	Panel   *Panel   `json:"-"`
	User    *User    `json:"-"`
}
