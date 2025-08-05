package services

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"strings"
)

// EmailService handles email operations
type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
	useSSL       bool
}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	useSSL := getEnvOrDefault("SMTP_USE_SSL", "true") == "true"
	defaultPort := "587"
	if useSSL {
		defaultPort = "465"
	}

	return &EmailService{
		smtpHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:     getEnvOrDefault("SMTP_PORT", defaultPort),
		smtpUsername: os.Getenv("SMTP_USERNAME"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromEmail:    getEnvOrDefault("FROM_EMAIL", os.Getenv("SMTP_USERNAME")),
		fromName:     getEnvOrDefault("FROM_NAME", "Tranza"),
		useSSL:       useSSL,
	}
}

// GenerateVerificationCode generates a 6-digit random verification code
func (es *EmailService) GenerateVerificationCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)
	
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		code[i] = digits[num.Int64()]
	}
	
	return string(code), nil
}

// SendVerificationEmail sends a verification code to the user's email
func (es *EmailService) SendVerificationEmail(to, username, code string) error {
	if es.smtpUsername == "" || es.smtpPassword == "" {
		// For development: log the code instead of sending email
		fmt.Printf("📧 [DEV MODE] Email verification code for %s (%s): %s\n", username, to, code)
		fmt.Printf("📧 [DEV MODE] Email would be sent from: %s\n", es.fromEmail)
		return nil
	}

	subject := "🔐 Your Tranza Verification Code"
	body := es.buildVerificationEmailBody(username, code)
	
	fmt.Printf("📧 Sending verification email to %s with code: %s\n", to, code)
	return es.sendEmail(to, subject, body)
}

// SendWelcomeEmail sends a welcome email after successful verification
func (es *EmailService) SendWelcomeEmail(to, username string) error {
	if es.smtpUsername == "" || es.smtpPassword == "" {
		// For development: log instead of sending email
		fmt.Printf("📧 Welcome email sent to %s (%s)\n", username, to)
		return nil
	}

	subject := "Welcome to Tranza!"
	body := es.buildWelcomeEmailBody(username)
	
	return es.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP with proper Gmail SSL/TLS support
func (es *EmailService) sendEmail(to, subject, body string) error {
	if es.smtpUsername == "" || es.smtpPassword == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	// Build the email message
	msg := es.buildEmailMessage(to, subject, body)
	addr := fmt.Sprintf("%s:%s", es.smtpHost, es.smtpPort)

	if es.useSSL && es.smtpPort == "465" {
		// For Gmail port 465 (SSL)
		return es.sendEmailSSL(addr, to, msg)
	} else {
		// For port 587 (STARTTLS)
		return es.sendEmailSTARTTLS(addr, to, msg)
	}
}

// sendEmailSSL sends email using SSL (port 465)
func (es *EmailService) sendEmailSSL(addr, to, msg string) error {
	// Create TLS connection
	tlsConfig := &tls.Config{
		ServerName: es.smtpHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, es.smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	auth := smtp.PlainAuth("", es.smtpUsername, es.smtpPassword, es.smtpHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set sender and recipient
	if err := client.Mail(es.fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send DATA command: %w", err)
	}

	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close message: %w", err)
	}

	return nil
}

// sendEmailSTARTTLS sends email using STARTTLS (port 587)
func (es *EmailService) sendEmailSTARTTLS(addr, to, msg string) error {
	auth := smtp.PlainAuth("", es.smtpUsername, es.smtpPassword, es.smtpHost)
	err := smtp.SendMail(addr, auth, es.fromEmail, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email via STARTTLS: %w", err)
	}
	return nil
}

// buildEmailMessage builds the email message with headers
func (es *EmailService) buildEmailMessage(to, subject, body string) string {
	var msg strings.Builder
	
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", es.fromName, es.fromEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)
	
	return msg.String()
}

// buildVerificationEmailBody creates the HTML body for verification email
func (es *EmailService) buildVerificationEmailBody(username, code string) string {
	expiryMinutes := getEnvOrDefault("EMAIL_VERIFICATION_EXPIRY_MINUTES", "15")
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification - Tranza</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6; 
            color: #333; 
            background-color: #f8fafc;
        }
        .container { 
            max-width: 600px; 
            margin: 40px auto; 
            background: white;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white; 
            padding: 40px 30px; 
            text-align: center; 
        }
        .header h1 { 
            font-size: 28px; 
            font-weight: 600; 
            margin-bottom: 8px;
        }
        .header p { 
            opacity: 0.9; 
            font-size: 16px;
        }
        .content { 
            padding: 40px 30px; 
        }
        .greeting { 
            font-size: 20px; 
            font-weight: 600; 
            margin-bottom: 20px;
            color: #1a202c;
        }
        .message { 
            font-size: 16px; 
            margin-bottom: 30px; 
            color: #4a5568;
        }
        .code-container {
            text-align: center;
            margin: 30px 0;
        }
        .code { 
            font-size: 36px; 
            font-weight: 700; 
            color: #667eea;
            letter-spacing: 8px; 
            margin: 20px 0; 
            padding: 20px 30px;
            background: linear-gradient(135deg, #f7fafc 0%%, #edf2f7 100%%);
            border: 2px solid #e2e8f0;
            border-radius: 12px;
            display: inline-block;
            min-width: 280px;
        }
        .instructions {
            background-color: #f7fafc;
            border-left: 4px solid #667eea;
            padding: 20px;
            margin: 25px 0;
            border-radius: 0 8px 8px 0;
        }
        .instructions h3 {
            color: #2d3748;
            margin-bottom: 10px;
            font-size: 16px;
        }
        .instructions p {
            color: #4a5568;
            margin-bottom: 8px;
        }
        .warning { 
            background-color: #fed7d7;
            border: 1px solid #feb2b2;
            color: #742a2a; 
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .warning strong {
            display: block;
            margin-bottom: 8px;
        }
        .warning ul {
            margin-left: 18px;
        }
        .warning li {
            margin-bottom: 4px;
        }
        .footer { 
            background-color: #f7fafc;
            text-align: center; 
            color: #718096; 
            font-size: 14px; 
            padding: 25px 30px;
            border-top: 1px solid #e2e8f0;
        }
        .brand {
            color: #667eea;
            font-weight: 600;
        }
        @media (max-width: 600px) {
            .container { margin: 20px; }
            .header, .content, .footer { padding: 25px 20px; }
            .code { font-size: 30px; letter-spacing: 6px; min-width: 240px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Verification Code</h1>
            <p>Secure your Tranza account</p>
        </div>
        <div class="content">
            <div class="greeting">Hi %s! 👋</div>
            
            <div class="message">
                Welcome to <strong class="brand">Tranza</strong>! To complete your registration and secure your account, please verify your email address using the verification code below:
            </div>
            
            <div class="code-container">
                <div class="code">%s</div>
            </div>
            
            <div class="instructions">
                <h3>📝 How to use this code:</h3>
                <p>1. Return to the Tranza registration page</p>
                <p>2. Enter the 6-digit code exactly as shown above</p>
                <p>3. Click "Verify Email" to activate your account</p>
            </div>
            
            <div class="warning">
                <strong>⚠️ Important Security Information:</strong>
                <ul>
                    <li>This code expires in <strong>%s minutes</strong></li>
                    <li>Never share this code with anyone</li>
                    <li>Tranza will never ask for this code via phone or chat</li>
                    <li>If you didn't request this code, please ignore this email</li>
                </ul>
            </div>
        </div>
        <div class="footer">
            <p>This is an automated security message from <strong class="brand">Tranza</strong></p>
            <p>Please do not reply to this email • Need help? Contact our support team</p>
        </div>
    </div>
</body>
</html>`, username, code, expiryMinutes)
}

// buildWelcomeEmailBody creates the HTML body for welcome email
func (es *EmailService) buildWelcomeEmailBody(username string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to Tranza</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #10b981; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f9fafb; padding: 30px; }
        .footer { text-align: center; color: #6b7280; font-size: 14px; margin-top: 20px; }
        .cta { background-color: #3b82f6; color: white; padding: 12px 24px; 
               text-decoration: none; border-radius: 5px; display: inline-block; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome to Tranza!</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>Congratulations! Your email has been verified and your account is now active.</p>
            
            <p>You can now:</p>
            <ul>
                <li>✅ Access your secure dashboard</li>
                <li>✅ Manage your financial transactions</li>
                <li>✅ Use all Tranza features</li>
            </ul>
            
            <p>Thank you for choosing Tranza for your financial needs!</p>
        </div>
        <div class="footer">
            <p>Welcome aboard! 🚀</p>
        </div>
    </div>
</body>
</html>`, username)
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}