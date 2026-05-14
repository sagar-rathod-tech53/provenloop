package utils

import (
	"crypto/tls"

	"github.com/sagar-rathod-tech53/provenloop/config"
	"gopkg.in/gomail.v2"
)

func SendEmail(
	to string,
	subject string,
	body string,
) error {

	m := gomail.NewMessage()

	m.SetHeader("From", config.AppConfig.EmailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		config.AppConfig.SMTPHost,
		config.AppConfig.SMTPPort,
		config.AppConfig.SMTPUser,
		config.AppConfig.SMTPPass,
	)

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return d.DialAndSend(m)
}
