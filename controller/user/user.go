package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

// Login 用户登录
func Login(c *gin.Context) {

}

//Settings 个人设置
func Settings(c *gin.Context) {
	type SettingForm struct {
		Name   string `binding:"required,min=2,max=12"`
		Phone  string `binding:"required,min=2,max=20"`
		Weixin string `binding:"required,min=2,max=20"`
		QQ     string `binding:"required,min=2,max=20"`
	}
	var lf SettingForm
	if err := c.ShouldBind(&lf); err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, "您的输入不符合规范，请检查后重试")
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)
	u.Weixin = lf.Weixin
	u.QQ = lf.QQ
	u.Phone = lf.Phone
	u.Name = lf.Name
	var err error
	if c.Request.Method == http.MethodPost {
		err = nbdomain.DB.Save(&u).Error
	} else {
		err = nbdomain.DB.Model(&u).Update(u).Error
	}
	if err != nil {
		log.Println("database error", err.Error())
		c.String(http.StatusInternalServerError, "服务器错误：数据库错误。")
		return
	}
	c.JSON(http.StatusOK, u)
}
