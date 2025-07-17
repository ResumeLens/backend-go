package utils

import (
	"fmt"
	"net/smtp"

	"github.com/spf13/viper"
)

// SendInviteEmail sends an invite link via email
func SendInviteEmail(recipientEmail, inviteToken string) error {
	smtpHost := viper.GetString("SMTP_HOST")
	smtpPort := viper.GetString("SMTP_PORT")
	smtpUser := viper.GetString("SMTP_USER")
	smtpPass := viper.GetString("SMTP_PASS")
	senderName := viper.GetString("SMTP_SENDER_NAME")

	from := smtpUser
	to := []string{recipientEmail}
	subject := "You're Invited to Join ResumeLens"
	inviteLink := fmt.Sprintf("https://resumelens.com/accept-invite?token=%s", inviteToken)

	body := fmt.Sprintf("Hello,\n\nYou've been invited to join ResumeLens.\n\nAccept your invite here: %s\n\nThis invite expires in 48 hours.\n\nBest,\n%s", inviteLink, senderName)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return smtp.SendMail(addr, auth, from, to, message)
}
