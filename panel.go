package panel

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
	gorm.Model
	UserID       uint
	Domain       string `gorm:"type:varchar(100);unique_index"`
	Name         string
	NameEn       string
	Desc         string
	DescEn       string
	AnalysisType string
	Analysis     string
	Theme        string
	OfferTheme   string

	User    User     `json:",omitempty"`
	Cats    []Cat    `json:",omitempty"`
	Domains []Domain `json:",omitempty"`
}

//SID 字符串ID
func (p Panel) SID() string {
	return fmt.Sprintf("%d", p.ID)
}
