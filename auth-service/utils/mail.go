package utils

import (
	"fmt"
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	from := GetEnv("SMTP_EMAIL", "")
	password := GetEnv("SMTP_PASSWORD", "")
	smtpHost := GetEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := GetEnv("SMTP_PORT", "587")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte("Subject: " + subject + "\r\n" +
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}
