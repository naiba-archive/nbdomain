package panel

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

//Offers 报价列表
func Offers(c *gin.Context) {
	u := c.MustGet(mygin.KUser).(model.User)
	nbdomain.DB.Model(&u).Related(&u.Offers)
	c.JSON(http.StatusOK, u.Offers)
}

//Export 导出米表
func Export(c *gin.Context) {
	u := c.MustGet(mygin.KUser).(model.User)
	var p model.Panel
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, c.Param("id")).First(&p).Error != nil {
		c.String(http.StatusForbidden, "米表不存在")
		return
	}
	nbdomain.DB.Model(&p).Related(&p.Cats)
	var txt = make([]byte, 0)
	for i := 0; i < len(p.Cats); i++ {
		nbdomain.DB.Model(&p.Cats[i]).Related(&p.Cats[i].Domains)
		txt = append(txt, []byte("#"+p.Cats[i].Name+","+p.Cats[i].NameEn+"\n")...)
		for j := 0; j < len(p.Cats[i].Domains); j++ {
			txt = append(txt, []byte(p.Cats[i].Domains[j].Domain+
				","+strconv.Itoa(p.Cats[i].Domains[j].Cost)+
				","+p.Cats[i].Domains[j].Buy.Format("2006-01-02")+
				","+strconv.Itoa(p.Cats[i].Domains[j].Renew)+
				","+p.Cats[i].Domains[j].Desc+"\n")...)
		}
	}
	c.String(http.StatusOK, string(txt))
}

//ListCats 米表分类列表
func ListCats(c *gin.Context) {
	u := c.MustGet(mygin.KUser).(model.User)
	var p model.Panel
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, c.Param("id")).First(&p).Error != nil {
		c.String(http.StatusForbidden, "米表不存在")
		return
	}
	nbdomain.DB.Model(&p).Related(&p.Cats)
	c.JSON(http.StatusOK, p.Cats)
}

//ListDomains 米表域名列表
func ListDomains(c *gin.Context) {
	u := c.MustGet(mygin.KUser).(model.User)
	var p model.Panel
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, c.Param("id")).First(&p).Error != nil {
		c.String(http.StatusForbidden, "米表不存在")
		return
	}
	nbdomain.DB.Model(&p).Related(&p.Domains)
	c.JSON(http.StatusOK, p.Domains)
}

//List 米表列表
func List(c *gin.Context) {
	u := c.MustGet(mygin.KUser).(model.User)
	nbdomain.DB.Model(&u).Related(&u.Panels)
	c.JSON(http.StatusOK, u.Panels)
}

//Delete 删除米表
func Delete(c *gin.Context) {
	id := c.Param("id")
	var p model.Panel
	u := c.MustGet(mygin.KUser).(model.User)
	if nbdomain.DB.Where("id = ? AND user_id = ?", id, u.ID).First(&p).Error != nil {
		c.String(http.StatusForbidden, "米表不存在")
		return
	}
	os.Remove("data/upload/logo/" + p.SID() + "-logo.png")
	os.Remove("data/upload/logo/" + p.SID() + "-logo_en.png")
	nbdomain.DB.Delete(&p)
}

//Edit 添加/修改米表
func Edit(c *gin.Context) {
	type PanelForm struct {
		ID           uint   `form:"id"`
		Domain       string `form:"domain" binding:"required,min=3,max=63"`
		Name         string `form:"name_cn" binding:"required,min=1,max=20"`
		NameEn       string `form:"name_en" binding:"required,min=1,max=40"`
		Desc         string `form:"desc_cn" binding:"required,min=1,max=255"`
		DescEn       string `form:"desc_en" binding:"required,min=1,max=1000"`
		Theme        string `form:"theme" binding:"required"`
		OfferTheme   string `form:"offer_theme" binding:"required"`
		Analysis     string `form:"ga" binding:"max=20"`
		AnalysisType string `form:"at"`
	}
	var pf PanelForm
	if e := c.ShouldBind(&pf); e != nil {
		log.Println(e)
		c.String(http.StatusForbidden, "输入数据不符合规范。")
		return
	}
	if pf.AnalysisType != "" {
		if _, has := model.AnalysisTypes[pf.AnalysisType]; !has {
			c.String(http.StatusForbidden, "米表统计类型不存在")
			return
		}
	}
	if _, has := model.ThemeList[pf.Theme]; !has {
		c.String(http.StatusForbidden, "主题不存在")
		return
	}
	if len(pf.Domain) < 4 {
		c.String(http.StatusForbidden, "域名格式不符合规范")
		return
	}

	u := c.MustGet(mygin.KUser).(model.User)

	// 查询会员是否有效
	if u.GoldVIPExpire.Before(time.Now()) && u.SuperVIPExpire.Before(time.Now()) {
		c.String(http.StatusForbidden, "您还不是会员，无法进行此操作")
		return
	}

	// 根据会员等级限制米表数量
	var panelCount int
	nbdomain.DB.Where("user_id = ?").Find(model.Panel{}).Count(&panelCount)
	if u.SuperVIPExpire.After(time.Now()) {
		// 限制数量
		if panelCount > 5 {
			c.String(http.StatusForbidden, "您的米表数超过5，无法进行此操作")
			return
		}
	} else {
		// 限制米表主题
		if pf.Theme != "offical-superhero" {
			c.String(http.StatusForbidden, "您是黄金会员，只能使用「superhero」主题")
			return
		}
		// 限制Offer主题
		if pf.OfferTheme != "offical-superhero" {
			c.String(http.StatusForbidden, "您是黄金会员，只能使用「superhero」主题")
			return
		}
		// 限制数量
		if panelCount > 1 {
			c.String(http.StatusForbidden, "您的米表数超过1，无法进行此操作，建议您升级会员")
			return
		}
	}

	//如果是修改米表，鉴权
	var p model.Panel
	if c.Request.Method == http.MethodPut {
		if nbdomain.DB.Where("id = ? AND user_id = ?", pf.ID, u.ID).First(&p).Error != nil {
			c.String(http.StatusForbidden, "米表不存在")
			return
		}
	}

	// 保存logo
	saveLogo := func(f *multipart.FileHeader, what string) error {
		ext := f.Filename[strings.LastIndex(f.Filename, ".")+1:]
		return c.SaveUploadedFile(f, "data/upload/logo/"+p.SID()+"-"+what+"."+ext)
	}

	// 检查logo
	checkLogo := func(name string) (*multipart.FileHeader, error, bool) {
		f, err := c.FormFile(name)
		if err != nil {
			if c.Request.Method == http.MethodPut {
				return f, nil, false
			}
			return f, errors.New("上传文件获取失败。"), false
		}
		fname := f.Filename
		if !strings.HasSuffix(fname, ".png") || f.Size > 1024*700 {
			return f, errors.New("logo只能是 png 格式，且不得超过700K"), false
		}
		return f, nil, true
	}
	flogo, e, blogo := checkLogo("logo_cn")
	if !blogo && e != nil {
		c.String(http.StatusForbidden, e.Error())
		return
	}
	flogoEn, e, blogoEn := checkLogo("logo_en")
	if !blogo && e != nil {
		c.String(http.StatusForbidden, e.Error())
		return
	}
	// 数据入库
	if c.Request.Method == http.MethodPost {
		p.UserID = u.ID
	}
	p.Name = pf.Name
	p.NameEn = pf.NameEn
	p.Domain = strings.ToLower(pf.Domain)
	p.Desc = pf.Desc
	p.DescEn = pf.DescEn
	p.Theme = pf.Theme
	p.OfferTheme = pf.OfferTheme
	p.Analysis = pf.Analysis
	p.AnalysisType = pf.AnalysisType
	var err error
	if c.Request.Method == http.MethodPost {
		err = nbdomain.DB.Save(&p).Error
	} else {
		err = nbdomain.DB.Model(&p).Update(p).Error
	}
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if blogo {
		saveLogo(flogo, "logo")
	}
	if blogoEn {
		saveLogo(flogoEn, "logo_en")
	}
	c.JSON(http.StatusOK, p)
}
