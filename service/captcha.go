package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/naiba/domain-panel"
)

//CaptchaService 验证码服务
type CaptchaService struct{}

type recaptchaResp struct {
	Success  bool
	Hostname string
}

//Verify 验证验证码
func (cs CaptchaService) Verify(gresp, ip string) (flag bool, host string) {
	resp, err := http.Post("https://www.recaptcha.net/recaptcha/api/siteverify",
		"application/x-www-form-urlencoded",
		strings.NewReader("secret="+panel.CF.ReCaptcha+"&response="+gresp+"&remoteip="+ip))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var rp recaptchaResp
	err = json.Unmarshal(body, &rp)
	if err != nil {
		return
	}
	return rp.Success, rp.Hostname
}
