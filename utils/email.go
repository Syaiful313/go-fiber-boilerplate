package utils

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

type EmailData struct {
	To      string
	Subject string
	Body    string
}

func SendEmail(config EmailConfig, emailData EmailData) error {
	// Setup authentication
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)

	// Compose the email
	to := []string{emailData.To}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		emailData.To, emailData.Subject, emailData.Body))

	// Send the email
	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.FromEmail, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func GenerateResetPasswordEmail(resetLink string) string {
	return fmt.Sprintf(`
		<html>
		<body>
			<h2>Reset Your Password</h2>
			<p>You have requested to reset your password. Click the link below to reset your password:</p>
			<p><a href="%s" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Reset Password</a></p>
			<p>If you did not request this, please ignore this email.</p>
			<p>This link will expire in 1 hour.</p>
		</body>
		</html>
	`, resetLink)
}

func GeneratePasswordResetSuccessEmail() string {
	return `
		<html>
		<body>
			<h2>Password Reset Successful</h2>
			<p>Your password has been successfully reset.</p>
			<p>If you did not perform this action, please contact our support team immediately.</p>
		</body>
		</html>
	`
}

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

