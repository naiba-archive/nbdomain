package domain

import (
	"log"
	"net/http"
	"time"

	panel "git.cm/nb/domain-panel"
	"git.cm/nb/domain-panel/pkg/mygin"
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
	if !panel.DomainRegexp.Match([]byte(ef.Domain)) {
		c.String(http.StatusForbidden, "域名格式不符合规范")
		return
	}
	u := c.MustGet(mygin.KUser).(panel.User)
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
	d.Buy = ef.Buy
	d.Cost = ef.Cost
	d.Renew = ef.Renew
	d.Registrar = ef.Registrar

	if err := panel.DB.Save(&d).Error; err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, d)
}
