package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("static", "theme/static")
	r.LoadHTMLGlob("theme/template/**/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "offical-superhero/index.html", nil)
	})
	r.Run()
}
