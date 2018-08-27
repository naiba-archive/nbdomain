package service

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"

	"git.cm/nb/domain-panel"
)

//MailService 邮件服务
type MailService struct{}

//SendMail 发送邮件
func (ms MailService) SendMail(toMail, subj, body string) (flag bool) {
	from := mail.Address{Address: panel.CF.Mail.User}
	to := mail.Address{Address: toMail}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	host, _, _ := net.SplitHostPort(panel.CF.Mail.SMTP)

	auth := smtp.PlainAuth("", panel.CF.Mail.User, panel.CF.Mail.Pass, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	c, err := smtp.Dial(panel.CF.Mail.SMTP)
	if err != nil {
		log.Println(err)
		return
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Println(err)
		return
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Println(err)
		return
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Println(err)
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println(err)
		return
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
		return
	}

	c.Quit()
	return true
}
