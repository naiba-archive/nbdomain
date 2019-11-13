package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/smartwalle/alipay"
	"golang.org/x/oauth2"

	"github.com/naiba/com"
	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/model"
	"github.com/naiba/nbdomain/pkg/mygin"
	"github.com/naiba/nbdomain/service"
)

func procNotify(nti *alipay.TradeNotification) error {
	var o model.Order
	if err := nbdomain.DB.Where("id = ?", nti.OutTradeNo).First(&o).Error; err != nil {
		return err
	}
	if !o.Finish {
		nbdomain.DB.Model(&o).Related(&o.User)
		if o.What == "gold" {
			if o.User.GoldVIPExpire.Before(time.Now()) {
				o.User.GoldVIPExpire = time.Now().Add(time.Hour * 24 * 30)
			} else {
				o.User.GoldVIPExpire = o.User.GoldVIPExpire.Add(time.Hour * 24 * 30)
			}
		} else {
			if o.User.SuperVIPExpire.Before(time.Now()) {
				o.User.SuperVIPExpire = time.Now().Add(time.Hour * 24 * 30)
			} else {
				o.User.SuperVIPExpire = o.User.SuperVIPExpire.Add(time.Hour * 24 * 30)
			}
		}
		return nbdomain.DB.Save(&o.User).Error
	}
	return nil
}

var oidcConfig *oidc.Config
var config oauth2.Config
var ctx context.Context
var verifier *oidc.IDTokenVerifier

func init() {
	ctx = context.Background()
	clientID := "1-66hM9Z"
	clientSecret := "twC7ItQTM31wqjJf"

	provider, err := oidc.NewProvider(ctx, "https://space.mentuo.com")
	if err != nil {
		log.Fatal(err)
	}
	oidcConfig = &oidc.Config{
		ClientID: clientID,
	}
	verifier = provider.Verifier(oidcConfig)
	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "https://" + nbdomain.CF.Web.Domain + "/hack/oauth2-redirect",
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}
}

// Oauth2Login 烧饼社群登录
func Oauth2Login(c *gin.Context) {
	state := com.RandomString(7)
	service.CacheService{}.Instance().Add("st-"+state, 1, time.Minute*5)
	c.Redirect(http.StatusFound, config.AuthCodeURL(state))
}

// Oauth2LoginCallback 烧饼社群登录回调
func Oauth2LoginCallback(c *gin.Context) {
	_, has := service.CacheService{}.Instance().Get("st-" + c.Query("state"))
	if !has {
		c.String(http.StatusBadRequest, "state did not match")
		return
	}

	oauth2Token, err := config.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to exchange token: "+err.Error())
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.String(http.StatusInternalServerError, "No id_token field in oauth2 token.")
		return
	}
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to verify ID Token: "+err.Error())
		return
	}

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var x map[string]map[string]interface{}
	err = json.Unmarshal(data, &x)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	oid := x["IDTokenClaims"]["sub"]
	username := x["IDTokenClaims"]["name"]
	var u model.User
	var newUser = false
	err = nbdomain.DB.Model(model.User{}).Where("ucenter_id = ?", oid).First(&u).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.String(http.StatusForbidden, err.Error())
			return
		}
		newUser = true
		u = model.User{
			UcenterID: oid.(string),
		}
	}
	if u.UcenterExtra != username.(string) {
		u.UcenterExtra = username.(string)
	}
	if newUser {
		nbdomain.DB.Create(&u)
		u.GoldVIPExpire = time.Now()
		u.SuperVIPExpire = u.GoldVIPExpire
		if u.ID == 1 {
			u.IsAdmin = true
		}
	}
	if err := u.GenerateToken(nbdomain.DB); err != nil {
		c.String(http.StatusInternalServerError, "数据库错误")
		return
	}
	c.JSON(http.StatusOK, u)
}

//Settings 个人设置
func Settings(c *gin.Context) {
	type SettingForm struct {
		Name   string `binding:"required,min=2,max=12"`
		Phone  string `binding:"required,min=2,max=20"`
		Weixin string `binding:"required,min=2,max=20"`
		QQ     string `binding:"required,min=2,max=20"`
	}
	var lf SettingForm
	if err := c.ShouldBind(&lf); err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, "您的输入不符合规范，请检查后重试")
		return
	}
	u := c.MustGet(mygin.KUser).(model.User)
	u.Weixin = lf.Weixin
	u.QQ = lf.QQ
	u.Phone = lf.Phone
	u.Name = lf.Name
	var err error
	if c.Request.Method == http.MethodPost {
		err = nbdomain.DB.Save(&u).Error
	} else {
		err = nbdomain.DB.Model(&u).Update(u).Error
	}
	if err != nil {
		log.Println("database error", err.Error())
		c.String(http.StatusInternalServerError, "服务器错误：数据库错误。")
		return
	}
	c.JSON(http.StatusOK, u)
}
