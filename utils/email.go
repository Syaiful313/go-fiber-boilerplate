package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
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
	if !ValidateEmail(emailData.To) || !ValidateEmail(config.FromEmail) {
		return fmt.Errorf("invalid email address")
	}

	toAddr, _ := mail.ParseAddress(emailData.To)
	fromAddr, _ := mail.ParseAddress(config.FromEmail)

	if err := rejectHeaderInjection(emailData.Subject); err != nil {
		return err
	}

	// Setup authentication
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)

	// Compose the email
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		toAddr.Address, emailData.Subject, emailData.Body))

	addr := net.JoinHostPort(config.SMTPHost, config.SMTPPort)
	tlsConfig := &tls.Config{
		ServerName: config.SMTPHost,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to establish TLS connection: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %v", err)
	}

	if err := client.Mail(fromAddr.Address); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err := client.Rcpt(toAddr.Address); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data command: %v", err)
	}

	if _, err := writer.Write(msg); err != nil {
		return fmt.Errorf("failed to write email body: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to finalize email: %v", err)
	}

	return client.Quit()
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
	if strings.ContainsAny(email, "\r\n") {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func rejectHeaderInjection(value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("invalid header value")
	}
	return nil
}
