package panel

import (
	"errors"
	"fmt"
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

type batchForm struct {
	PanelID uint64 `binding:"required,min=1" json:"panel_id,omitempty"`
	Cats    []struct {
		Name    string `binding:"required,min=1,max=20" json:"name,omitempty"`
		NameEn  string `binding:"required,min=1,max=30" json:"name_en,omitempty"`
		Domains []struct {
			Cost   int       `binding:"min=1" json:"cost,omitempty"` //购入成本
			Buy    time.Time `json:"buy,omitempty"`
			Renew  int       `binding:"min=1" json:"renew,omitempty"` //续费成本
			Domain string    `binding:"required,min=1,max=64" json:"domain,omitempty"`
			Desc   string    `binding:"required,min=1,max=200" json:"desc,omitempty"`
		} `json:"domains,omitempty"`
	} `json:"cats,omitempty"`
}

//Import 批量导入域名
func Import(c *gin.Context) {
	var bf batchForm
	if err := c.ShouldBind(&bf); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("数据不符合规范：%s", err.Error()),
		})
		return
	}
	for _, cat := range bf.Cats {
		for _, domain := range cat.Domains {
			if len(domain.Domain) < 4 {
				c.JSON(http.StatusOK, model.Response{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("数据不符合规范：%s", domain.Domain),
				})
				return
			}
		}
	}
	u := c.MustGet(mygin.KUser).(model.User)
	var p model.Panel
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, bf.PanelID).First(&p).Error != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("米表不存在：%d", bf.PanelID),
		})
		return
	}
	addedDomains := make([]model.Domain, 0)
	for _, catForm := range bf.Cats {
		var cat model.Cat
		if nbdomain.DB.Where("name = ? AND user_id = ?", strings.TrimSpace(catForm.Name), u.ID).First(&cat).Error != nil {
			cat.Name = catForm.Name
			cat.NameEn = catForm.NameEn
			cat.UserID = u.ID
			cat.PanelID = p.ID
			if err := nbdomain.DB.Save(&cat).Error; err != nil {
				c.JSON(http.StatusOK, model.Response{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("数据库错误：%s", err),
				})
				return
			}
		}
		for _, domainForm := range catForm.Domains {
			var domain model.Domain
			if nbdomain.DB.Where("domain = ?", domainForm.Domain).First(&domain).Error == nil {
				continue
			}
			domain.UserID = u.ID
			domain.PanelID = p.ID
			domain.CatID = cat.ID
			domain.Buy = &domainForm.Buy
			domain.Cost = domainForm.Cost
			domain.Renew = domainForm.Renew
			domain.Domain = domainForm.Domain
			domain.Desc = domainForm.Desc
			if err := nbdomain.DB.Save(&domain).Error; err != nil {
				c.JSON(http.StatusOK, model.Response{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("数据库错误：%s", err),
				})
				return
			}
			addedDomains = append(addedDomains, domain)
		}
	}
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: len(addedDomains),
	})
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

type listPanelReq struct {
	UserID uint64 `form:"-"`
	Domain string `form:"domain"`
	Name   string `form:"name"`
}

//List 米表列表
func List(c *gin.Context) {
	var lpr listPanelReq
	if err := c.ShouldBindQuery(&lpr); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)
	lpr.UserID = u.ID
	var respList model.ListData
	var ts []model.Panel
	if err := model.AfterPagination(model.WhereQuery(nbdomain.DB.Model(model.Panel{}), lpr), model.BeforePagenation(c), &respList).Find(&ts).Error; err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	for i := 0; i < len(ts); i++ {
		ts[i].Stat(nbdomain.DB)
	}
	respList.List = ts
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: respList,
	})
}

//Delete 删除米表
func Delete(c *gin.Context) {
	id := c.Param("id")
	var p model.Panel
	u := c.MustGet(mygin.KUser).(model.User)
	if nbdomain.DB.Where("id = ? AND user_id = ?", id, u.ID).First(&p).Error != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "米表不存在",
		})
		return
	}
	os.Remove("data/upload/logo/" + p.SID() + "-logo.png")
	os.Remove("data/upload/logo/" + p.SID() + "-logo_en.png")
	nbdomain.DB.Delete(&model.Domain{
		PanelID: p.ID,
	})
	nbdomain.DB.Delete(&model.Cat{
		PanelID: p.ID,
	})
	if e := nbdomain.DB.Delete(&p).Error; e != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("数据库错误：%s", e.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: id,
	})
}

type panelForm struct {
	ID           uint64 `form:"id"`
	Domain       string `form:"domain" binding:"required,min=3,max=63"`
	Name         string `form:"name" binding:"required,min=1,max=20"`
	NameEn       string `form:"name_en" binding:"required,min=1,max=40"`
	Desc         string `form:"desc" binding:"required,min=1,max=255"`
	DescEn       string `form:"desc_en" binding:"required,min=1,max=1000"`
	Theme        string `form:"theme" binding:"required"`
	OfferTheme   string `form:"offer_theme" binding:"required"`
	Analysis     string `form:"analysis" binding:"max=20"`
	AnalysisType string `form:"analysis_type"`
}

//Edit 添加/修改米表
func Edit(c *gin.Context) {
	var pf panelForm
	if e := c.ShouldBind(&pf); e != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("输入数据不符合规范：%s", e.Error()),
		})
		return
	}
	if pf.AnalysisType != "" {
		if _, has := model.AnalysisTypes[pf.AnalysisType]; !has {
			c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusBadRequest,
				Message: "米表统计类型不存在",
			})
			return
		}
	}
	if _, has := model.ThemeList[pf.Theme]; !has {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "主题不存在",
		})
		return
	}
	if len(pf.Domain) < 4 {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "域名格式不符合规范",
		})
		return
	}

	u := c.MustGet(mygin.KUser).(model.User)

	//如果是修改米表，鉴权
	var p model.Panel
	if pf.ID != 0 {
		if nbdomain.DB.Where("id = ? AND user_id = ?", pf.ID, u.ID).First(&p).Error != nil {
			c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusBadRequest,
				Message: "米表不存在",
			})
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
			if pf.ID != 0 {
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

	flogo, e, blogo := checkLogo("logo")
	if !blogo && e != nil && pf.ID == 0 {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("文件上传失败：%s", e.Error()),
		})
		return
	}

	flogoEn, e, blogoEn := checkLogo("logo_en")
	if !blogoEn && e != nil && pf.ID == 0 {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("文件上传失败：%s", e.Error()),
		})
		return
	}
	// 数据入库
	if pf.ID == 0 {
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
	if pf.ID == 0 {
		err = nbdomain.DB.Save(&p).Error
	} else {
		err = nbdomain.DB.Model(&p).Update(p).Error
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("数据库错误：%s", err.Error()),
		})
		return
	}
	if blogo {
		saveLogo(flogo, "logo")
	}
	if blogoEn {
		saveLogo(flogoEn, "logo_en")
	}

	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: p.ID,
	})
}
