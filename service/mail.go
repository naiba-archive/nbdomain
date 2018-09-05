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
	from := mail.Address{Name: "lifelongsender", Address: panel.CF.Mail.User}
	to := mail.Address{Address: toMail}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj + "（通知邮件）"
	headers["Content-Type"] = "text/plain; charset=utf-8"

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

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", panel.CF.Mail.SMTP, tlsconfig)
	if err != nil {
		log.Println(err)
		return
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Println(err)
		return
	}

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
