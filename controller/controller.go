package controller

import (
	"html/template"
	"net/http"
	"strings"

	"git.cm/nb/domain-panel"
	"git.cm/nb/domain-panel/controller/cat"
	"git.cm/nb/domain-panel/controller/domain"
	"git.cm/nb/domain-panel/controller/mibiao"
	"git.cm/nb/domain-panel/controller/panelr"
	"git.cm/nb/domain-panel/controller/user"
	"git.cm/nb/domain-panel/controller/verify"
	"git.cm/nb/domain-panel/controller/whois"
	"git.cm/nb/domain-panel/pkg/mygin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Web start
func Web() {
	var mode string
	if panel.CF.Debug {
		mode = gin.DebugMode
	} else {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"toLower": strings.ToLower,
	})
	r.LoadHTMLGlob("theme/template/**/*")
	if panel.CF.Debug {
		conf := cors.DefaultConfig()
		conf.AllowAllOrigins = true
		conf.AddAllowMethods("DELETE")
		conf.AddAllowHeaders("Authorization")
		r.Use(cors.New(conf))
		r.Static("static", "theme/static")
		r.Static("upload", "upload/")
	}
	panelRouter := r.Group("/")
	{
		panelRouter.GET("", mibiao.Index)
		panelRouter.GET("offer/:domain", mibiao.Offer)
		panelRouter.POST("offer/:domain", mibiao.Offer)
		panelRouter.GET("offer/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/")
		})
		panelRouter.GET("allowed", mibiao.Allow)
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
			authUser.PUT("user", user.Settings)
			authUser.GET("offers", panelr.Offers)
			authUser.PUT("panel", panelr.Edit)
			authUser.GET("themes", func(c *gin.Context) {
				c.JSON(http.StatusOK, panel.ThemeList)
			})
			authUser.GET("analysis_types", func(c *gin.Context) {
				c.JSON(http.StatusOK, panel.AnalysisTypes)
			})
			authUser.DELETE("panel/:id", panelr.Delete)
			authUser.GET("panel/:id/cats", panelr.ListCats)
			authUser.GET("panel/:id/domains", panelr.ListDomains)
			authUser.POST("panel", panelr.Edit)
			authUser.DELETE("cat/:id", cat.Delete)
			authUser.PUT("cat", cat.Edit)
			authUser.POST("cat", cat.Edit)
			authUser.GET("panels", panelr.List)
			authUser.POST("batch", domain.Batch)
			authUser.POST("domain", domain.Edit)
			authUser.PUT("domain", domain.Edit)
			authUser.DELETE("domain/:id", domain.Delete)
			authUser.GET("whois/:domain", whois.Whois)
		}
	}
	go r.Run(panel.CF.Web.Addr)
}
