package user

import (
	"log"
	"net/http"

	"git.cm/nb/domain-panel/pkg/mygin"

	"git.cm/nb/domain-panel"
	"git.cm/nb/domain-panel/service"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

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
	u := c.MustGet(mygin.KUser).(panel.User)
	u.Weixin = lf.Weixin
	u.QQ = lf.QQ
	u.Phone = lf.Phone
	u.Name = lf.Name
	if err := panel.DB.Save(&u).Error; err != nil {
		log.Println("database error", err.Error())
		c.String(http.StatusInternalServerError, "服务器错误：数据库错误。")
		return
	}
	c.JSON(http.StatusOK, u)
}

//Login 登录
func Login(ctx *gin.Context) {
	type loginForm struct {
		Mail      string `form:"mail" binding:"required,email"`
		Password  string `form:"password" binding:"required,min=6"`
		ReCaptcha string `form:"recaptcha" binding:"required,min=20"`
	}
	var lf loginForm
	if err := ctx.ShouldBind(&lf); err != nil {
		log.Println(err)
		ctx.String(http.StatusForbidden, "您的输入不符合规范，请检查后重试")
		return
	}
	var cs service.CaptchaService
	if success, host := cs.Verify(lf.ReCaptcha, ctx.ClientIP()); !success || host != panel.CF.Web.Domain {
		ctx.String(http.StatusForbidden, "验证码不正确")
		return
	}
	var u panel.User
	if panel.DB.Where("mail = ?", lf.Mail).First(&u).Error != nil {
		ctx.String(http.StatusForbidden, "用户不存在")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(lf.Password)) != nil {
		ctx.String(http.StatusForbidden, "密码不正确")
		return
	}
	if err := u.GenerateToken(); err != nil {
		log.Println("database error", err.Error())
		ctx.String(http.StatusInternalServerError, "服务器错误：数据库错误。")
		return
	}
	ctx.JSON(http.StatusOK, u)
}

//Register 注册账号
func Register(ctx *gin.Context) {
	type regForm struct {
		Name     string `form:"name" binding:"required,min=2,max=12"`
		Mail     string `form:"mail" binding:"required,email"`
		Password string `form:"password" binding:"required,min=6"`
		Verify   string `form:"verify" binding:"required,len=5"`
	}
	var rf regForm
	if err := ctx.ShouldBind(&rf); err != nil {
		log.Println(err)
		ctx.String(http.StatusForbidden, "您的输入不符合规范，请检查后重试")
		return
	}
	//校验验证码
	cacheKey := "v" + "reg" + rf.Mail + rf.Verify
	var cs service.CacheService
	if _, has := cs.Instance().Get(cacheKey); !has {
		ctx.String(http.StatusForbidden, "邮箱验证码不正确")
		return
	}
	cs.Instance().Delete(cacheKey)
	//用户入库
	var u panel.User
	u.Name = rf.Name
	u.Mail = rf.Mail
	bPass, err := bcrypt.GenerateFromPassword([]byte(rf.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("password generate", err.Error())
		ctx.String(http.StatusInternalServerError, "服务器错误：密码生成错误。")
		return
	}
	u.Password = string(bPass)
	if err := panel.DB.Save(&u).Error; err != nil {
		log.Println("database error", err.Error())
		ctx.String(http.StatusInternalServerError, "服务器错误：数据库错误。")
		return
	}
	if u.ID == 1 {
		u.IsAdmin = true
	}
	if err := u.GenerateToken(); err != nil {
		log.Println("database error", err.Error())
		ctx.String(http.StatusInternalServerError, "服务器错误：数据库错误。")
		return
	}
	ctx.JSON(http.StatusOK, u)
}

//ResetPassword 重置密码
func ResetPassword(ctx *gin.Context) {
	type resetForm struct {
		Mail     string `form:"mail" binding:"required,email"`
		Password string `form:"password" binding:"required,min=6"`
		Verify   string `form:"verify" binding:"required,len=5"`
	}
	var rf resetForm
	if err := ctx.ShouldBind(&rf); err != nil {
		log.Println(err)
		ctx.String(http.StatusForbidden, "您的输入不符合规范，请检查后重试")
		return
	}
	//校验验证码
	cacheKey := "v" + "forget" + rf.Mail + rf.Verify
	var cs service.CacheService
	if _, has := cs.Instance().Get(cacheKey); !has {
		ctx.String(http.StatusForbidden, "邮箱验证码不正确")
		return
	}
	cs.Instance().Delete(cacheKey)
	//用户入库
	var u panel.User
	if panel.DB.Where("mail = ?", rf.Mail).First(&u).Error != nil {
		ctx.String(http.StatusForbidden, "用户不存在")
		return
	}
	bPass, err := bcrypt.GenerateFromPassword([]byte(rf.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("password generate", err.Error())
		ctx.String(http.StatusInternalServerError, "服务器错误：密码生成错误。")
		return
	}
	u.Password = string(bPass)
	if err := u.GenerateToken(); err != nil {
		log.Println("database error", err.Error())
		ctx.String(http.StatusInternalServerError, "服务器错误：数据库错误。")
		return
	}
	ctx.JSON(http.StatusOK, u)
}
