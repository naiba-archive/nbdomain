package cat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

type loadCatReq struct {
	UserID  uint64 `form:"-"`
	PanelID uint64 `form:"panel_id"`
	Name    string `form:"name"`
}

//List 分类列表
func List(c *gin.Context) {
	var lcr loadCatReq
	if err := c.ShouldBindQuery(&lcr); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)
	lcr.UserID = u.ID
	var respList model.ListData
	var ts []model.Cat
	if err := model.AfterPagination(model.WhereQuery(nbdomain.DB.Model(model.Cat{}),
		lcr), model.BeforePagenation(c), &respList).Find(&ts).Error; err != nil {
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

//Delete 删除分类
func Delete(c *gin.Context) {
	id := c.Param("id")
	u := c.MustGet(mygin.KUser).(model.User)
	var cat model.Cat
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, id).First(&cat).Error != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "分类不存在。",
		})
		return
	}
	nbdomain.DB.Where("cat_id = ?", id).Delete(&model.Domain{})
	nbdomain.DB.Delete(&cat)
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: id,
	})
}

type editForm struct {
	PanelID uint64 `json:"panel_id" binding:"required"`
	ID      uint64 `json:"id"`
	Index   int    `json:"index"`
	Name    string `json:"name" binding:"required,min=1,max=20"`
	NameEn  string `json:"name_en" binding:"required,min=1,max=30"`
}

//Edit 修改或添加分类
func Edit(c *gin.Context) {
	var ef editForm
	if err := c.ShouldBind(&ef); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("输入数据不符合规范：%s", err.Error()),
		})
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)
	var p model.Panel
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, ef.PanelID).First(&p).Error != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "米表不存在。",
		})
		return
	}
	var cat model.Cat
	if ef.ID != 0 {
		if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, ef.ID).First(&cat).Error != nil {
			c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusBadRequest,
				Message: "分类不存在。",
			})
			return
		}
	}
	cat.Name = ef.Name
	cat.NameEn = ef.NameEn
	cat.UserID = u.ID
	cat.PanelID = p.ID
	cat.Index = ef.Index
	var err error
	if c.Request.Method == http.MethodPost {
		err = nbdomain.DB.Save(&cat).Error
	} else {
		err = nbdomain.DB.Model(&cat).Update(cat).Error
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("服务器错误：%s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: cat.ID,
	})
}
