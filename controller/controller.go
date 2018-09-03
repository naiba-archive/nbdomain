package controller

import (
	"git.cm/nb/domain-panel"
	"git.cm/nb/domain-panel/controller/panelr"
	"git.cm/nb/domain-panel/controller/user"
	"git.cm/nb/domain-panel/controller/verify"
	"git.cm/nb/domain-panel/pkg/mygin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Web start
func Web() {
	r := gin.Default()
	r.LoadHTMLGlob("theme/template/**/*")
	if panel.CF.Debug {
		conf := cors.DefaultConfig()
		conf.AllowAllOrigins = true
		conf.AddAllowMethods("DELETE")
		conf.AddAllowHeaders("Authorization")
		r.Use(cors.New(conf))
		//用户上传文件
		r.Static("/upload", "upload")
	}
	api := r.Group("/api")
	{
		unAuth := api.Group("")
		{
			unAuth.Use(mygin.Authorize(mygin.AuthOption{NeedGhost: true}))
			unAuth.POST("send_verify", verify.Mail)
			unAuth.POST("reg", user.Register)
			unAuth.POST("login", user.Login)
			unAuth.POST("reset_password", user.ResetPassword)
		}
		authUser := api.Group("")
		{
			authUser.Use(mygin.Authorize(mygin.AuthOption{NeedUser: true}))
			authUser.PUT("panel", panelr.Edit)
			authUser.DELETE("panel/:id", panelr.Delete)
			authUser.POST("panel", panelr.Edit)
			authUser.GET("panels", panelr.List)
		}
	}
	go r.Run(panel.CF.Web.Addr)
}
