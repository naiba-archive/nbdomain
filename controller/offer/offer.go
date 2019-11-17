package offer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

// Delete 删除Offer
func Delete(c *gin.Context) {
	id := c.Param("id")
	u := c.MustGet(mygin.KUser).(model.User)
	var d model.Offer
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, id).First(&d).Error != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: "Offer 不存在。",
		})
		return
	}
	nbdomain.DB.Delete(&d)
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: id,
	})
}

type loadOfferReq struct {
	UserID uint64 `form:"-"`
	Domain string `form:"domain"`
}

//List 域名列表
func List(c *gin.Context) {
	var ldr loadOfferReq
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
	var ts []model.Offer
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
