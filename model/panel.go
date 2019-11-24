package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//ThemeList 主题列表
var ThemeList = map[string]string{
	"offical-superhero": "「官方」Superhero",
	"sbdotsb":           "「烧饼」Doma.in",
	"longge":            "「龙哥」青山为雪白头",
	"bocai":             "「神马」接地气菠菜",
}

//OfferThemeList 报价页主题列表
var OfferThemeList = map[string]string{
	"offical-superhero": "「官方」Superhero",
	"sbdotsb":           "「烧饼」Doma.in",
}

//AnalysisTypes 米表统计类型
var AnalysisTypes = map[string]string{
	"g": "谷歌统计",
	"c": "CNZZ",
}

//Panel 米表
type Panel struct {
	Common
	UserID       uint64 `json:"user_id,omitempty"`
	Domain       string `gorm:"type:varchar(100);unique_index" json:"domain,omitempty"`
	Name         string `json:"name,omitempty"`
	NameEn       string `json:"name_en,omitempty"`
	Desc         string `json:"desc,omitempty"`
	DescEn       string `json:"desc_en,omitempty"`
	AnalysisType string `json:"analysis_type,omitempty"`
	Analysis     string `json:"analysis,omitempty"`
	Theme        string `json:"theme,omitempty"`
	OfferTheme   string `json:"offer_theme,omitempty"`

	TotalRenew uint64   `gorm:"-" json:"total_renew,omitempty"`
	TotalBuy   uint64   `gorm:"-" json:"total_buy,omitempty"`
	User       User     `json:"-"`
	Cats       []Cat    `json:"-"`
	Domains    []Domain `json:"-"`
}

type sumResult struct {
	Trenew uint64
	Tbuy   uint64
}

// Stat ..
func (p *Panel) Stat(db *gorm.DB) {
	var r sumResult
	db.Model(Domain{}).Select("sum(cost) as tbuy,sum(renew) as trenew").Where("panel_id = ?", p.ID).Scan(&r)
	p.TotalRenew = r.Trenew
	p.TotalBuy = r.Tbuy
}

//SID 字符串ID
func (p Panel) SID() string {
	return fmt.Sprintf("%d", p.ID)
}
