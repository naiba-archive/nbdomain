package mygin

import (
	panel "git.cm/nb/domain-panel"
	"github.com/gin-gonic/gin"
)

//SetCookie 设置Cookie
func SetCookie(k, v string, ctx *gin.Context) {
	ctx.SetCookie(k, v, 60*60*24*365*2, "/", panel.CF.Web.Domain, false, false)
}
