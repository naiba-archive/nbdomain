package controller

import (
	"net/http"

	"git.cm/nb/domain-panel"

	"git.cm/nb/domain-panel/controller/user"
	"github.com/gin-gonic/gin"
)

//Web start
func Web() {
	r := gin.Default()
	r.Static("static", "theme/static")
	r.LoadHTMLGlob("theme/template/**/*")
	unAuth := r.Group("")
	{
		unAuth.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "offical/index.html", nil)
		})
		unAuth.POST("/login", user.Login)
	}
	go r.Run(panel.CF.Web.Addr)
}
