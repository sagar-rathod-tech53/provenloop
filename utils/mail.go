package utils

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/sagar-rathod-devops/do-host-network-backend/config"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// ? Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(to, subject, body string) error {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return err
	}

	from := cfg.EmailFrom
	smtpPass := cfg.SMTPPass
	smtpUser := cfg.SMTPUser
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	// Create new message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Dialer configuration
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // NOTE: for dev only

	// Send email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully to", to)
	return nil
}

// GenerateOTP generates a random OTP of the given length.
func GenerateOTP(length int) string {
	rand.Seed(int64(time.Now().UnixNano()))
	otp := ""
	for i := 0; i < length; i++ {
		otp += fmt.Sprintf("%d", rand.Intn(10))
	}
	return otp
}
