package controller

import (
	"git.cm/nb/domain-panel"
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
	r.Use(cors.Default())
	api := r.Group("/api")
	{
		unAuth := api.Group("")
		unAuth.Use(mygin.Authorize(mygin.AuthOption{NeedGhost: true}))
		unAuth.POST("send_verify", verify.Mail)
		unAuth.POST("reg", user.Register)
		unAuth.POST("login", user.Login)
		unAuth.POST("reset_password", user.ResetPassword)
	}
	go r.Run(panel.CF.Web.Addr)
}
