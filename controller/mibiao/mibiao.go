package mibiao

import (
	"net/http"
	"strings"
	"time"

	"git.cm/nb/domain-panel"

	"github.com/gin-gonic/gin"
)

func checkExpire(c *gin.Context) bool {
	domain := stripPort(c.Request.Host)
	var p panel.Panel
	err := panel.DB.Where("domain = ?", domain).First(&p).Error
	if err != nil {
		// 域名
		var d panel.Domain
		err = panel.DB.Where("domain = ?", domain).First(&d).Error
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "https://"+panel.CF.Web.Domain)
			return false
		}
		panel.DB.Model(&d).Related(&d.User)
		if d.User.Expire.Before(time.Now()) {
			c.String(http.StatusOK, "域名停放已过期，请您及时续费。")
			return false
		}
		panel.DB.Model(&d).Related(&d.Panel)
		c.Redirect(http.StatusTemporaryRedirect, "https://"+d.Panel.Domain+"/offer/"+domain)
		return false
	}
	panel.DB.Model(&p).Related(&p.User)
	if p.User.Expire.Before(time.Now()) {
		c.String(http.StatusOK, "域名停放已过期，请您及时续费。")
		return false
	}
	c.Set("Panel", p)
	c.Set("Chinese", strings.Contains(c.Request.Header.Get("accept-language"), "zh"))
	return true
}

//Allow 米表自动HTTPS
func Allow(c *gin.Context) {
	if c.ClientIP() != "127.0.0.1" {
		c.Status(http.StatusForbidden)
		return
	}
	domain := c.Query("domain")
	var p panel.Panel
	err := panel.DB.Where("domain = ?", domain).First(&p).Error
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	panel.DB.Model(&p).Related(&p.User)
	if p.User.Expire.Before(time.Now()) {
		c.Status(http.StatusForbidden)
		return
	}
}

//Index 米表首页
func Index(c *gin.Context) {
	if !checkExpire(c) {
		return
	}
	p := c.MustGet("Panel").(panel.Panel)
	panel.DB.Model(&p).Related(&p.Cats)
	for i := 0; i < len(p.Cats); i++ {
		panel.DB.Model(&p.Cats[i]).Related(&p.Cats[i].Domains)
	}
	isChinese := c.GetBool("Chinese")
	var title string
	if isChinese {
		title = p.Name
	} else {
		title = p.NameEn
	}
	c.HTML(http.StatusOK, "offical-superhero/index", gin.H{
		"Title":   title,
		"Panel":   p,
		"Chinese": isChinese,
	})
}

//Offer 报价页
func Offer(c *gin.Context) {
	if !checkExpire(c) {
		return
	}
	p := c.MustGet("Panel").(panel.Panel)
	var d panel.Domain
	if panel.DB.Where("domain = ?", c.Param("domain")).First(&d).Error != nil {
		c.Redirect(http.StatusTemporaryRedirect, "https://"+p.Domain)
		return
	}
	isChinese := c.GetBool("Chinese")
	var title string
	if isChinese {
		title = d.Domain + " 或许可以出售 - " + p.NameEn
	} else {
		title = d.Domain + " may be for sale - " + p.NameEn
	}
	c.HTML(http.StatusOK, "offical-superhero/offer", gin.H{
		"Title":   title,
		"Panel":   p,
		"Domain":  d,
		"Chinese": isChinese,
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
