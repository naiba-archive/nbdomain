package controller

import (
	"git.cm/nb/domain-panel"

	"github.com/gin-gonic/gin"
)

//Web start
func Web() {
	r := gin.Default()
	r.LoadHTMLGlob("theme/template/**/*")

	go r.Run(panel.CF.Web.Addr)
}
