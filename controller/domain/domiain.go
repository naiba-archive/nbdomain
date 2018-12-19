package domain

import (
	"log"
	"net/http"
	"strings"
	"time"

	panel "git.cm/naiba/domain-panel"
	"git.cm/naiba/domain-panel/pkg/mygin"
	"github.com/gin-gonic/gin"
)

//Delete 删除域名
func Delete(c *gin.Context) {
	id := c.Param("id")
	u := c.MustGet(mygin.KUser).(panel.User)
	var d panel.Domain
	if panel.DB.Where("user_id = ? AND id = ?", u.ID, id).First(&d).Error != nil {
		c.String(http.StatusForbidden, "域名不存在。")
		return
	}
	panel.DB.Delete(&d)
}

//Batch 批量导入域名
func Batch(c *gin.Context) {
	type BatchForm struct {
		PanelID uint `binding:"required,min=1"`
		Cats    []struct {
			Name    string `binding:"required,min=1,max=20"`
			NameEn  string `binding:"required,min=1,max=30"`
			Domains []struct {
				Cost   int `binding:"min=1"` //购入成本
				Buy    time.Time
				Renew  int    `binding:"min=1"` //续费成本
				Domain string `binding:"required,min=1,max=64"`
				Desc   string `binding:"required,min=1,max=200"`
			}
		}
	}
	var bf BatchForm
	if err := c.ShouldBind(&bf); err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, "输入数据不符合规范。可留空但不可以乱填。")
		return
	}
	for _, cat := range bf.Cats {
		for _, domain := range cat.Domains {
			if len(domain.Domain) < 4 {
				c.String(http.StatusForbidden, domain.Domain+":域名格式不符合规范")
				return
			}
		}
	}
	u := c.MustGet(mygin.KUser).(panel.User)
	var p panel.Panel
	if panel.DB.Where("user_id = ? AND id = ?", u.ID, bf.PanelID).First(&p).Error != nil {
		c.String(http.StatusForbidden, "米表不存在。")
		return
	}
	addedDomains := make([]panel.Domain, 0)
	for _, catForm := range bf.Cats {
		var cat panel.Cat
		if panel.DB.Where("name = ? AND user_id = ?", strings.TrimSpace(catForm.Name), u.ID).First(&cat).Error != nil {
			cat.Name = catForm.Name
			cat.NameEn = catForm.NameEn
			cat.UserID = u.ID
			cat.PanelID = p.ID
			if panel.DB.Save(&cat).Error != nil {
				c.String(http.StatusInternalServerError, "数据库错误，联系管理员")
				return
			}
		}
		for _, domainForm := range catForm.Domains {
			var domain panel.Domain
			if panel.DB.Where("domain = ?", domainForm.Domain).First(&domain).Error == nil {
				continue
			}
			domain.UserID = u.ID
			domain.PanelID = p.ID
			domain.CatID = cat.ID
			domain.Buy = domainForm.Buy
			domain.Cost = domainForm.Cost
			domain.Renew = domainForm.Renew
			domain.Domain = domainForm.Domain
			domain.Desc = domainForm.Desc
			if panel.DB.Save(&domain).Error != nil {
				c.String(http.StatusInternalServerError, "数据库错误，联系管理员")
				return
			}
			addedDomains = append(addedDomains, domain)
		}
	}
	c.JSON(http.StatusOK, addedDomains)
}

//Edit 添加/修改域名
func Edit(c *gin.Context) {
	type EditForm struct {
		CatID     uint `binding:"required,min=1"`
		PanelID   uint `binding:"required,min=1"`
		ID        uint
		Create    time.Time //注册时间
		Expire    time.Time //到期时间
		Cost      int       `binding:"min=1"` //购入成本
		Renew     int       `binding:"min=1"` //续费成本
		Buy       time.Time //购入时间
		Registrar string    `binding:"min=1,max=100"`
		Domain    string    `binding:"required,min=1,max=64"`
		Desc      string    `binding:"required,min=1,max=200"`
	}
	var ef EditForm
	if err := c.ShouldBind(&ef); err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, "输入数据不符合规范。")
		return
	}
	if len(ef.Domain) < 4 {
		c.String(http.StatusForbidden, "域名格式不符合规范")
		return
	}
	u := c.MustGet(mygin.KUser).(panel.User)

	// 查询会员是否有效
	if u.GoldVIPExpire.Before(time.Now()) && u.SuperVIPExpire.Before(time.Now()) {
		c.String(http.StatusForbidden, "您还不是会员，无法进行此操作")
		return
	}

	// 根据会员等级限制域名数量
	var domainCount int
	panel.DB.Where("user_id = ?").Find(panel.Domain{}).Count(&domainCount)
	if u.SuperVIPExpire.After(time.Now()) {
		if domainCount > 1000 {
			c.String(http.StatusForbidden, "您的域名数超过1000，无法进行此操作")
			return
		}
	} else {
		if domainCount > 100 {
			c.String(http.StatusForbidden, "您的域名数超过100，无法进行此操作，建议您升级会员")
			return
		}
	}

	var cat panel.Cat
	if panel.DB.Where("user_id = ? AND id = ?", u.ID, ef.CatID).First(&cat).Error != nil {
		c.String(http.StatusForbidden, "分类不存在。")
		return
	}

	var d panel.Domain
	if ef.ID != 0 {
		if panel.DB.Where("user_id = ? AND id = ?", u.ID, ef.ID).First(&d).Error != nil {
			c.String(http.StatusForbidden, "域名不存在。")
			return
		}
	}

	d.CatID = ef.CatID
	d.PanelID = ef.PanelID
	d.Domain = ef.Domain
	d.Desc = ef.Desc
	d.UserID = u.ID
	d.Create = ef.Create
	d.Expire = ef.Expire
	if d.Expire.After(time.Now()) {
		d.WhoisUpdate = time.Now()
	}
	d.Buy = ef.Buy
	d.Cost = ef.Cost
	d.Renew = ef.Renew
	d.Registrar = ef.Registrar

	var err error
	if c.Request.Method == http.MethodPost {
		err = panel.DB.Save(&d).Error
	} else {
		err = panel.DB.Model(&d).Update(d).Error
	}
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, d)
}
