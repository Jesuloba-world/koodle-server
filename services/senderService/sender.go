package senderservice

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/k3a/html2text"
	"github.com/mailersend/mailersend-go"

	"github.com/Jesuloba-world/koodle-server/model"
	userRepo "github.com/Jesuloba-world/koodle-server/repo/user"
)

type SenderService struct {
	client   *mailersend.Mailersend
	sender   string
	userRepo *userRepo.UserRepo
}

func NewSenderService(apiKey, sender string, userRepo *userRepo.UserRepo) *SenderService {
	client := mailersend.NewMailersend(apiKey)
	return &SenderService{
		client:   client,
		sender:   sender,
		userRepo: userRepo,
	}
}

func (e *SenderService) sendEmail(to, subject, htmlBody string) error {
	plainTextBody := html2text.HTML2Text(htmlBody)

	recipients := []mailersend.Recipient{
		{
			Email: to,
			Name:  "Recipient",
		},
	}

	from := mailersend.From{
		Email: e.sender,
		Name:  "Koodle",
	}

	message := e.client.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetHTML(htmlBody)
	message.SetText(plainTextBody)
	message.SetSubject(subject)

	_, err := e.client.Email.Send(context.Background(), message)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *SenderService) generateHTMLContent(templateFilePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template :%w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (e *SenderService) SendOTPToRecipient(channel model.OTPChannel, purpose model.OTPPurpose, recipient, otp string, expirationDuration time.Duration) error {
	// otpFormatted := fmt.Sprintf("%s-%s", otp[:2], otp[2:])

	user, err := e.userRepo.FindByID(context.Background(), recipient)
	if err != nil {
		return fmt.Errorf("OTP recipient could not be found: %v", err)
	}

	switch channel {
	case model.OTPChannelEmail:
		fmt.Printf("Sending OTP %s to %s via email\n", otp, user.Email)

		var subject, templateFile string
		data := map[string]interface{}{
			"OTP":            otp,
			"ExpirationTime": expirationDuration.Minutes(),
		}

		switch purpose {
		case model.OTPPurposeEmailVerification:
			subject = "Verify your email"
			templateFile = "email_verification.html"
		case model.OTPPurposePasswordReset:
			subject = "Password reset"
			templateFile = "password_reset.html"
		}

		workingDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		templateFilePath := filepath.Join(workingDir, "templates/emailTemplates", templateFile)

		htmlContent, err := e.generateHTMLContent(templateFilePath, data)
		if err != nil {
			return err
		}

		err = e.sendEmail(user.Email, subject, htmlContent)
		if err != nil {
			return fmt.Errorf("failed to send OTP email: %w", err)
		}

	case model.OTPChannelSMS:
		// TODO: implement send sms logic if needed
		fmt.Printf("Sending OTP %s to %s via SMS\n", otp, recipient)
	}

	return nil
}
