package user

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

// GET ..
func GET(c *gin.Context) {
	u := c.MustGet(mygin.KUser).(model.User)
	var r model.Response
	r.Code = http.StatusOK
	r.Result = u
	c.JSON(http.StatusOK, r)
}

type loginForm struct {
	Mail      string `json:"mail" binding:"required|email"`
	Password  string `json:"password" binding:"required"`
	AutoLogin bool   `json:"auto_login"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var r model.Response

	var lf loginForm
	if err := c.ShouldBindJSON(&lf); err != nil {
		r.Code = http.StatusBadRequest
		r.Message = fmt.Sprintf("数据填写有误：%s", err)
		c.JSON(http.StatusOK, r)
		return
	}

	var u model.User
	if err := nbdomain.DB.Where("id = ?", 1).First(&u).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			r.Code = http.StatusInternalServerError
			r.Message = fmt.Sprintf("数据库出错啦：%s", err)
			c.JSON(http.StatusOK, r)
			return
		}
		// 第一个登录用户 管理员
		u.ID = 1
		u.Mail = lf.Mail
		pass, err := bcrypt.GenerateFromPassword([]byte(lf.Password), bcrypt.DefaultCost)
		if err != nil {
			r.Code = http.StatusInternalServerError
			r.Message = fmt.Sprintf("密码生成错误：%s", err)
			c.JSON(http.StatusOK, r)
			return
		}
		u.Password = string(pass)
		if err = nbdomain.DB.Create(&u).Error; err != nil {
			r.Code = http.StatusInternalServerError
			r.Message = fmt.Sprintf("数据库出错啦：%s", err)
			c.JSON(http.StatusOK, r)
			return
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(lf.Password))
	if err != nil || lf.Mail == "" || u.Mail != lf.Mail {
		r.Code = http.StatusBadRequest
		r.Message = "邮箱或密码错误"
		c.JSON(http.StatusOK, r)
		return
	}

	// 生成 Token
	if err = u.GenerateToken(nbdomain.DB, lf.AutoLogin); err != nil {
		r.Code = http.StatusInternalServerError
		r.Message = fmt.Sprintf("数据库出错啦：%s", err)
		c.JSON(http.StatusOK, r)
		return
	}

	r.Code = http.StatusOK
	r.Result = u
	c.JSON(http.StatusOK, r)
}

type settingForm struct {
	Name     string `binding:"omitempty,min=2,max=12" json:"name,omitempty"`
	Phone    string `binding:"omitempty,min=2,max=20" json:"phone,omitempty"`
	Weixin   string `binding:"omitempty,min=2,max=20" json:"weixin,omitempty"`
	QQ       string `binding:"omitempty,min=2,max=20" json:"qq,omitempty"`
	Password string `binding:"omitempty,min=6" json:"password,omitempty"`
}

//Settings 个人设置
func Settings(c *gin.Context) {
	var lf settingForm
	if err := c.ShouldBind(&lf); err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("输入数据不符合规范：%s", err.Error()),
		})
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)
	u.Weixin = lf.Weixin
	u.QQ = lf.QQ
	u.Phone = lf.Phone
	u.Name = lf.Name
	if lf.Password != "" {
		pass, err := bcrypt.GenerateFromPassword([]byte(lf.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("密码生成错误：%s", err),
			})
			return
		}
		u.Password = string(pass)
	}
	var err error
	err = nbdomain.DB.Model(&u).Update(u).Error
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("服务器错误：%s", err),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Code:   http.StatusOK,
		Result: u.ID,
	})
}
