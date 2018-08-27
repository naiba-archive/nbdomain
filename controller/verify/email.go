package verify

import (
	"net/http"

	"git.cm/nb/domain-panel/service"
	"github.com/gin-gonic/gin"
)

type mailForm struct {
	Type string `binding:"required,string"`
	Mail string `binding:"required,email"`
}

//Mail 邮箱验证码
func Mail(ctx *gin.Context) {
	var mf mailForm
	if err := ctx.ShouldBind(&mf); err != nil {
		ctx.String(http.StatusForbidden, "输入数据有误，请重试。")
		return
	}
	service.MailService.SendMail(mf.Mail, "")
}
