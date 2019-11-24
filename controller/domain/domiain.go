package domain

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

type loadDomainReq struct {
	UserID  uint64 `form:"-"`
	PanelID uint64 `form:"panel_id"`
	CatID   uint64 `form:"cat_id"`
	Domain  string `form:"domain"`
}

//List 域名列表
func List(c *gin.Context) {
	var ldr loadDomainReq
	if err := c.ShouldBindQuery(&ldr); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)
	ldr.UserID = u.ID
	var respList model.ListData
	var ts []model.Domain
	if err := model.AfterPagination(model.WhereQuery(nbdomain.DB.Model(model.Domain{}),
		ldr), model.BeforePagenation(c), &respList).Find(&ts).Error; err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	respList.List = ts
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: respList,
	})
}

//Delete 删除域名
func Delete(c *gin.Context) {
	id := c.Param("id")
	u := c.MustGet(mygin.KUser).(model.User)
	var d model.Domain
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, id).First(&d).Error != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "域名不存在。",
		})
		return
	}
	nbdomain.DB.Delete(&d)
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: id,
	})
}

type editForm struct {
	ID      uint64    `json:"id,omitempty"`
	CatID   uint64    `binding:"required,min=1" json:"cat_id,omitempty"`
	PanelID uint64    `binding:"required,min=1" json:"panel_id,omitempty"`
	Domain  string    `binding:"required,min=1,max=64" json:"domain,omitempty"`
	Desc    string    `binding:"required,min=1,max=200" json:"desc,omitempty"`
	Cost    int       `json:"cost,omitempty"`  //购入成本
	Renew   int       `json:"renew,omitempty"` //续费成本
	Buy     time.Time `json:"buy,omitempty"`   //购入时间

	Registrar string    `binding:"min=1,max=100" json:"registrar,omitempty"`
	Create    time.Time `json:"create,omitempty"` //注册时间
	Expire    time.Time `json:"expire,omitempty"` //到期时间
}

//Edit 添加/修改域名
func Edit(c *gin.Context) {
	var ef editForm
	if err := c.ShouldBind(&ef); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("输入数据不符合规范：%s", err.Error()),
		})
		return
	}
	if len(ef.Domain) < 4 {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "域名格式不符合规范",
		})
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)

	var cat model.Cat
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, ef.CatID).First(&cat).Error != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "分类不存在。",
		})
		return
	}

	var d model.Domain
	if ef.ID != 0 {
		if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, ef.ID).First(&d).Error != nil {
			c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusBadRequest,
				Message: "域名不存在。",
			})
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
	if ef.ID == 0 {
		err = nbdomain.DB.Save(&d).Error
	} else {
		err = nbdomain.DB.Model(&d).Update(d).Error
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("服务器错误：%s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: d.ID,
	})
}
