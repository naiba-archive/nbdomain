package mibiao

import (
	"net/http"
	"strings"
	"time"

	"git.cm/nb/domain-panel"

	"github.com/gin-gonic/gin"
)

func checkExpire(c *gin.Context) {
	domain := stripPort(c.Request.Host)
	var p panel.Panel
	err := panel.DB.Where("domain = ?", domain).First(&p).Error
	if err != nil {
		// 域名
		var d panel.Domain
		err = panel.DB.Where("domain = ?", domain).First(&d).Error
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "https://"+panel.CF.Web.Domain)
			return
		}
		panel.DB.Model(&d).Related(&d.User)
		panel.DB.Model(&d).Related(&d.Panel)
		if d.User.Expire.Before(time.Now()) {
			c.String(http.StatusOK, "您的会员服务已过期，请您及时续费。")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, "https://"+d.Panel.Domain+"/offer/"+domain)
		return
	}
	panel.DB.Model(&p).Related(&p.User)
	if p.User.Expire.Before(time.Now()) {
		c.String(http.StatusOK, "您的会员服务已过期，请您及时续费。")
		return
	}
	c.Set("Panel", p)
	c.Set("Chinese", strings.Contains(c.Request.Header.Get("accept-language"), "zh"))
}

//Index 米表首页
func Index(c *gin.Context) {
	checkExpire(c)
	p := c.MustGet("Panel").(panel.Panel)
	panel.DB.Model(&p).Related(&p.Cats)
	for i := 0; i < len(p.Cats); i++ {
		panel.DB.Model(&p.Cats[i]).Related(&p.Cats[i].Domains)
	}
	c.HTML(http.StatusOK, "offical-superhero/index", gin.H{
		"Panel":   p,
		"Chinese": c.GetBool("Chinese"),
	})
}

//Offer 报价页
func Offer(c *gin.Context) {
	checkExpire(c)
	p := c.MustGet("Panel").(panel.Panel)
	var d panel.Domain
	if panel.DB.Where("domain = ?", c.Param("domain")).First(&d).Error != nil {
		c.Redirect(http.StatusTemporaryRedirect, "https://"+p.Domain)
		return
	}
	c.HTML(http.StatusOK, "offical-superhero/offer", gin.H{
		"Panel":   p,
		"Domain":  d,
		"Chinese": c.GetBool("Chinese"),
	})
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}
