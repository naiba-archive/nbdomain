package cat

import (
	"log"
	"net/http"

	"git.cm/nb/domain-panel"
	"git.cm/nb/domain-panel/pkg/mygin"

	"github.com/gin-gonic/gin"
)

//Delete 删除分类
func Delete(c *gin.Context) {
	id := c.Param("id")
	u := c.MustGet(mygin.KUser).(panel.User)
	var cat panel.Cat
	if panel.DB.Where("user_id = ? AND id = ?", u.ID, id).First(&cat).Error != nil {
		c.String(http.StatusForbidden, "分类不存在。")
		return
	}
	panel.DB.Where("cat_id = ?", id).Delete(&panel.Domain{})
	panel.DB.Delete(&cat)
}

//Edit 修改或添加分类
func Edit(c *gin.Context) {
	type EditForm struct {
		PanelID uint   `json:"pid" binding:"required"`
		ID      uint   `json:"id"`
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
	u := c.MustGet(mygin.KUser).(panel.User)
	var p panel.Panel
	if panel.DB.Where("user_id = ? AND id = ?", u.ID, ef.PanelID).First(&p).Error != nil {
		c.String(http.StatusForbidden, "米表不存在。")
		return
	}
	var cat panel.Cat
	if ef.ID != 0 {
		if panel.DB.Where("user_id = ? AND id = ?", u.ID, ef.ID).First(&cat).Error != nil {
			c.String(http.StatusForbidden, "分类不存在。")
			return
		}
	}
	cat.Name = ef.Name
	cat.NameEn = ef.NameEn
	cat.UserID = u.ID
	cat.PanelID = p.ID
	cat.Index = ef.Index

	if err := panel.DB.Save(&cat).Error; err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "服务器错误")
		return
	}
	c.JSON(http.StatusOK, cat)
}
