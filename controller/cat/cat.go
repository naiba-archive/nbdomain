package cat

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

//Delete 删除分类
func Delete(c *gin.Context) {
	id := c.Param("id")
	u := c.MustGet(mygin.KUser).(model.User)
	var cat model.Cat
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, id).First(&cat).Error != nil {
		c.String(http.StatusForbidden, "分类不存在。")
		return
	}
	nbdomain.DB.Where("cat_id = ?", id).Delete(&model.Domain{})
	nbdomain.DB.Delete(&cat)
}

//Edit 修改或添加分类
func Edit(c *gin.Context) {
	type EditForm struct {
		PanelID uint64   `json:"pid" binding:"required"`
		ID      uint64   `json:"id"`
		Name    string `json:"name" binding:"required,min=1,max=20"`
		NameEn  string `json:"name_en" binding:"required,min=1,max=30"`
		Index   int    `json:"index"`
	}
	var ef EditForm
	if err := c.ShouldBind(&ef); err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, "输入数据不符合规范。")
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)
	var p model.Panel
	if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, ef.PanelID).First(&p).Error != nil {
		c.String(http.StatusForbidden, "米表不存在。")
		return
	}
	var cat model.Cat
	if ef.ID != 0 {
		if nbdomain.DB.Where("user_id = ? AND id = ?", u.ID, ef.ID).First(&cat).Error != nil {
			c.String(http.StatusForbidden, "分类不存在。")
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
		log.Println(err)
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, cat)
}
