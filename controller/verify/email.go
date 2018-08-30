package verify

import (
	"log"
	"net/http"
	"time"

	"git.cm/nb/domain-panel"

	"github.com/naiba/com"

	"git.cm/nb/domain-panel/service"
	"github.com/gin-gonic/gin"
)

type mailForm struct {
	Type      string `form:"type" binding:"required,alpha"`
	Mail      string `form:"mail" binding:"required,email"`
	ReCaptcha string `form:"recaptcha" binding:"required,min=20"`
}

var typeMessage = map[string]string{
	"reg":    "润措注册验证码",
	"forget": "润措重置密码",
}

//Mail 邮箱验证码
func Mail(ctx *gin.Context) {
	var mf mailForm
	if err := ctx.ShouldBind(&mf); err != nil {
		ctx.String(http.StatusForbidden, "输入数据不符合规范，请重试")
		return
	}
	var rs = service.CaptchaService{}
	if success, host := rs.Verify(mf.ReCaptcha, ctx.ClientIP()); !success || host != panel.CF.Web.Domain {
		log.Println("ReCaptcha Resp =====>", success, host, ctx.Request.URL.Host)
		ctx.String(http.StatusForbidden, "验证码错误，请重试")
		return
	}
	var msg string
	var has bool
	if msg, has = typeMessage[mf.Type]; !has {
		ctx.String(http.StatusForbidden, "不存在的验证码类型")
		return
	}
	//重置密码邮件
	if mf.Type == "forget" {
		var u panel.User
		if panel.DB.Where("mail = ?", mf.Mail).First(&u).Error != nil {
			ctx.String(http.StatusForbidden, "用户不存在")
			return
		}
	}
	code := com.RandomString(5)
	ms := service.MailService{}
	cs := service.CacheService{}
	if !ms.SendMail(mf.Mail, msg, "您的验证码为："+code) {
		ctx.String(http.StatusInternalServerError, "邮件发送失败，请联系客服处理")
		return
	}
	//邮箱验证码
	cs.Instance().Add("v"+mf.Type+mf.Mail+code, 0, time.Minute*10)
}
