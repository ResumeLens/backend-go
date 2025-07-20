package utils

import (
	"fmt"
	"net/smtp"

	"github.com/resumelens/authservice/internal/config"
)

func SendInviteEmail(recipientEmail, inviteToken string, cfg *config.Config) error {
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort
	smtpUser := cfg.SMTPUser
	smtpPass := cfg.SMTPPass
	senderName := cfg.SMTPSenderName

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
