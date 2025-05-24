package services

import (
	"auth-service/internal/config"
	"auth-service/internal/models"
	"fmt"
	"net/smtp"
)

const emailTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Код подтверждения</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f0f8ff;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background-image: linear-gradient(135deg, #f5f7fa 0%%, #c3cfe2 100%%);
        }
        .container {
            background-color: #fff;
            padding: 30px;
            border-radius: 12px;
            box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
            text-align: center;
            max-width: 400px;
            width: 100%%;
            border: 1px solid #e1e8ed;
        }
        h2 {
            color: #2c3e50;
            font-size: 24px;
            margin-bottom: 20px;
        }
        p {
            color: #7f8c8d;
            font-size: 16px;
            margin-bottom: 30px;
        }
        .code {
            font-size: 32px;
            font-weight: bold;
            color: #3498db;
            letter-spacing: 5px;
            margin: 20px 0;
        }
        .footer {
            margin-top: 20px;
            font-size: 14px;
            color: #bdc3c7;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>Ваш код подтверждения</h2>
        <div class="code">%s</div>
        <p>Код действителен в течение 5 минут.</p>
        <div class="footer">
            <p>Если вы не запрашивали этот код, пожалуйста, проигнорируйте это сообщение.</p>
        </div>
    </div>
</body>
</html>`

type EmailService struct {
	config config.SMTPConfig
}

func NewEmailService(config config.SMTPConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

func (s *EmailService) SendVerificationCode(to string, code string, verificationType models.VerificationType) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	subject := "Подтверждение "
	switch verificationType {
	case models.VerificationTypeRegistration:
		subject += "регистрации"
	case models.VerificationTypeLogin:
		subject += "входа"
	case models.VerificationTypePassword:
		subject += "смены пароля"
	}

	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		emailTemplate, to, subject, code)

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, []byte(message))
}
