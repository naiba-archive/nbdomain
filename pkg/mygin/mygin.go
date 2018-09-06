package mygin

import (
	"net/http"
	"strings"

	panel "git.cm/nb/domain-panel"
	"github.com/gin-gonic/gin"
)

//KChinese 语言
const KChinese = "auth-is-chinese"

//KUser 认证用户
const KUser = "auth-user"

//AuthOption 鉴权配置
type AuthOption struct {
	NeedAdmin bool
	NeedUser  bool
	NeedGhost bool
}

//Authorize 鉴权
func Authorize(o AuthOption) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var u panel.User
		token := ctx.Request.Header.Get("Authorization")
		if len(token) != 32 || panel.DB.Where("token = ?", token).First(&u).Error != nil {
			if o.NeedAdmin || o.NeedUser {
				ctx.String(http.StatusUnauthorized, "登录状态已失效，请您重新登录。")
				ctx.Abort()
				return
			}
		} else {
			//只能游客访问
			if o.NeedGhost {
				ctx.String(http.StatusForbidden, "此页面只允许游客访问")
				ctx.Abort()
				return
			}
			//只能管理员访问
			if o.NeedAdmin && !u.IsAdmin {
				ctx.String(http.StatusForbidden, "此页面只允许管理员访问")
				ctx.Abort()
				return
			}
			ctx.Set(KUser, u)
		}
		ctx.Set(KChinese, strings.Contains(ctx.Request.Header.Get("accept-language"), "zh"))
		ctx.Next()
	}
}

//TmplData 模板数据
type TmplData struct {
	Title string
	Data  map[string]interface{}
}

//TmplVars 模板变量
func TmplVars(ctx *gin.Context, data TmplData) gin.H {
	all := make(map[string]interface{})
	if len(data.Title) > 0 {
		all["title"] = data.Title + " - 日落域名资产 riluo.cn"
	} else {
		all["title"] = "日落域名资产管理平台 riluo.cn"
	}
	all["user"], _ = ctx.Get(KUser)
	all["chinese"], _ = ctx.Get(KChinese)
	all["data"] = data.Data
	return all
}

//SetCookie 设置Cookie
func SetCookie(k, v string, ctx *gin.Context) {
	panic("can't set cookie in API")
	//ctx.SetCookie(k, v, 60*60*24*365*2, "/", panel.CF.Web.Domain, false, false)
}
