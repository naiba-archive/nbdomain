package controller

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/controller/cat"
	"github.com/naiba/nbdomain/controller/domain"
	"github.com/naiba/nbdomain/controller/mibiao"
	"github.com/naiba/nbdomain/controller/panel"
	"github.com/naiba/nbdomain/controller/user"
	"github.com/naiba/nbdomain/controller/whois"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
)

// Web start
func Web() {
	var mode string
	if nbdomain.CF.Debug {
		mode = gin.DebugMode
	} else {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"toLower": strings.ToLower,
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},
	})
	r.LoadHTMLGlob("theme/template/**/*")
	r.Static("static", "theme/static")
	if nbdomain.CF.Debug {
		conf := cors.DefaultConfig()
		conf.AllowAllOrigins = true
		conf.AddAllowMethods("DELETE")
		conf.AddAllowHeaders("Authorization")
		r.Use(cors.New(conf))
		r.Static("upload", "data/upload")
	}
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "http://"+c.Request.URL.Hostname())
	})

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
		api.POST("login", user.Login)

		authUser := api.Group("")
		{
			authUser.Use(mygin.Authorize(mygin.AuthOption{NeedUser: true}))
			authUser.GET("user", user.GET)
			authUser.PUT("user", user.Settings)
			authUser.GET("offers", panel.Offers)
			authUser.PUT("panel", panel.Edit)
			authUser.GET("themes", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"themes":       model.ThemeList,
					"offer_themes": model.OfferThemeList,
				})
			})
			authUser.GET("analysis_types", func(c *gin.Context) {
				c.JSON(http.StatusOK, model.AnalysisTypes)
			})
			authUser.DELETE("panel/:id", panel.Delete)
			authUser.GET("panel/:id/cats", panel.ListCats)
			authUser.GET("panel/:id/domains", panel.ListDomains)
			authUser.GET("panel/:id/export", panel.Export)
			authUser.POST("panel", panel.Edit)
			authUser.DELETE("cat/:id", cat.Delete)
			authUser.PUT("cat", cat.Edit)
			authUser.POST("cat", cat.Edit)
			authUser.GET("panels", panel.List)
			authUser.POST("batch", domain.Batch)
			authUser.POST("domain", domain.Edit)
			authUser.PUT("domain", domain.Edit)
			authUser.DELETE("domain/:id", domain.Delete)
			authUser.GET("whois/:domain", whois.Whois)
		}
	}
	go r.Run(nbdomain.CF.Web.Addr)
}
