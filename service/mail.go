package service

import (
	"log"

	"git.cm/nb/domain-panel"

	"github.com/matcornic/hermes"
	"gopkg.in/gomail.v2"
)

//TextMail 纯文本邮件
const TextMail = 1

//HTMLMail 网页邮件
const HTMLMail = 2

var mailRender = hermes.Hermes{
	// Optional Theme
	// Theme: new(Default)
	Product: hermes.Product{
		// Appears in header & footer of e-mails
		Name: "日落域名资产管理平台",
		Link: "https://www.riluo.cn/",
		Logo: "https://www.riluo.cn/static/offical/logo.png",
	},
}

//MailService 邮件服务
type MailService struct{}

//SendMail 发送邮件
func (ms MailService) SendMail(toMail, subj string, mail hermes.Email, mType int) (flag bool) {
	var contentType, mailBody string
	var err error
	switch mType {
	case TextMail:
		contentType = "text/plain"
		mailBody, err = mailRender.GeneratePlainText(mail)
		if err != nil {
			log.Println(err)
			return false
		}
	case HTMLMail:
		contentType = "text/html"
		mailBody, err = mailRender.GenerateHTML(mail)
		if err != nil {
			log.Println(err)
			return false
		}
	default:
		return false
	}
	m := gomail.NewMessage()
	m.SetHeader("From", panel.CF.Mail.User)
	m.SetHeader("To", toMail)
	// m.SetAddressHeader("Cc", panel.CF.Mail.User, "LifelongSender")
	m.SetHeader("Subject", subj)
	m.SetBody(contentType, mailBody)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewPlainDialer(panel.CF.Mail.SMTP, panel.CF.Mail.Port, panel.CF.Mail.User, panel.CF.Mail.Pass)
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return false
	}
	return true
}
