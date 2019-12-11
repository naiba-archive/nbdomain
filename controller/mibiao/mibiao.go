package mibiao

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/idna"

	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/service"
)

func checkRedirect(c *gin.Context) bool {
	var err error
	domain := stripPort(c.Request.Host)
	if strings.Contains(domain, "xn--") {
		domain, err = idna.ToUnicode(domain)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "https://"+nbdomain.CF.Web.Domain)
			return false
		}
	}
	var p model.Panel
	err = nbdomain.DB.Where("domain = ?", domain).First(&p).Error
	topLevelDomain := domainutil.Domain(domain)
	if err != nil {
		if topLevelDomain != "" {
			err = nbdomain.DB.Where("domain = ?", topLevelDomain).First(&p).Error
		} else {
			topLevelDomain = domain
		}
	}
	if err != nil {
		//不是米表，试试域名
		var d model.Domain
		err = nbdomain.DB.Where("domain = ?", topLevelDomain).First(&d).Error
		//未找到域名，跳转平台首页
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "https://"+nbdomain.CF.Web.Domain)
			return false
		}
		//取米表详情
		nbdomain.DB.Model(&d).Related(&d.Panel)
		c.Redirect(http.StatusTemporaryRedirect, "https://"+d.Panel.Domain+"/offer/"+domain)
		return false
	}
	nbdomain.DB.Model(&p).Related(&p.User)
	//设置默认主题
	if p.Theme == "" {
		p.Theme = "offical-superhero"
	}
	if p.OfferTheme == "" {
		p.OfferTheme = "offical-superhero"
	}
	c.Set("Panel", p)
	c.Set("Chinese", strings.Contains(c.Request.Header.Get("accept-language"), "zh"))
	// 强制 HTTPS
	c.Writer.Header().Add("Strict-Transport-Security", fmt.Sprintf("max-age=%d; preload", 60*60*24*7*18))
	return true
}

//Allow 米表自动HTTPS
func Allow(c *gin.Context) {
	// 禁止局域网之外的访问
	// if c.ClientIP() != "127.0.0.1" {
	// 	c.Status(http.StatusForbidden)
	// 	return
	// }
	domain := c.Query("domain")
	var p model.Panel
	err := nbdomain.DB.Where("domain = ?", domain).First(&p).Error
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	nbdomain.DB.Model(&p).Related(&p.User)
}

//Index 米表首页
func Index(c *gin.Context) {
	if !checkRedirect(c) {
		return
	}
	p := c.MustGet("Panel").(model.Panel)
	nbdomain.DB.Model(&p).Order("index").Association("cats").Find(&p.Cats)
	for i := 0; i < len(p.Cats); i++ {
		nbdomain.DB.Model(&p.Cats[i]).Related(&p.Cats[i].Domains)
	}
	isChinese := c.GetBool("Chinese")
	var title string
	if isChinese {
		title = p.Name
	} else {
		title = p.NameEn
	}
	c.HTML(http.StatusOK, p.Theme+"/index", gin.H{
		"Static":  "https://" + nbdomain.CF.Web.Domain,
		"Title":   title,
		"Panel":   p,
		"Chinese": isChinese,
	})
}

//Offer 报价页
func Offer(c *gin.Context) {
	if !checkRedirect(c) {
		return
	}
	p := c.MustGet("Panel").(model.Panel)
	var d model.Domain
	if nbdomain.DB.Where("domain = ?", c.Param("domain")).First(&d).Error != nil {
		if c.Request.Method == http.MethodGet {
			c.Redirect(http.StatusTemporaryRedirect, "https://"+p.Domain)
		} else {
			c.String(http.StatusForbidden, "域名未找到")
		}
		return
	}
	if c.Request.Method == http.MethodGet {
		//展示报价页
		isChinese := c.GetBool("Chinese")
		var title string
		if isChinese {
			title = d.Domain + " 或许可以出售 - " + p.Name
		} else {
			title = d.Domain + " may be for sale - " + p.NameEn
		}
		c.HTML(http.StatusOK, p.OfferTheme+"/offer", gin.H{
			"Static":    "https://" + nbdomain.CF.Web.Domain,
			"Title":     title,
			"Panel":     p,
			"Domain":    d,
			"ReCaptcha": nbdomain.CF.ReCaptcha.Key,
			"Chinese":   isChinese,
		})
	} else {
		// 提交报价
		type OfferForm struct {
			Name      string `binding:"required,min=2,max=20"`
			Amount    int    `binding:"required,min=1"`
			Mail      string `binding:"required,email"`
			Currency  string `binding:"required,containsany=$¥€"`
			Recaptcha string `json:"g-recaptcha-response" form:"g-recaptcha-response" binding:"required"`
		}
		var of OfferForm
		if err := c.ShouldBind(&of); err != nil {
			log.Println(err)
			c.String(http.StatusForbidden, "信息填写不规范。")
			return
		}
		var cs service.CaptchaService
		if y, host := cs.Verify(of.Recaptcha, c.ClientIP()); !y || host != p.Domain {
			c.String(http.StatusForbidden, "ReCaptcha验证未通过。")
			return
		}
		var o model.Offer
		o.Amount = of.Amount
		o.Currency = of.Currency
		o.Domain = c.Param("domain")
		o.UserID = d.UserID
		o.Name = of.Name
		o.Mail = of.Mail
		if err := nbdomain.DB.Save(&o).Error; err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "服务器错误")
			return
		}
		c.Status(http.StatusOK)
	}
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
