package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, body string) error {
	emailFrom := os.Getenv("EMAIL_FROM")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	headers := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n", emailFrom, to, subject)
	msg := headers + body

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := smtpHost + ":587"

	if err := smtp.SendMail(addr, auth, emailFrom, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}