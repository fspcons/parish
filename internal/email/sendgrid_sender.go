package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DefaultSendGridAPIURL is the production SendGrid v3 mail send endpoint.
const DefaultSendGridAPIURL = "https://api.sendgrid.com/v3/mail/send"

// SendGridSender sends transactional mail via SendGrid’s HTTP API.
// SendGrid is commonly used on Google Cloud (API key from Secret Manager, outbound HTTPS from Cloud Run/GKE).
type SendGridSender struct {
	APIKey     string
	FromEmail  string
	HTTPClient *http.Client
	// SendURL is the SendGrid mail send endpoint (from config or tests).
	SendURL string
}

// NewSendGridSender returns a sender that uses SendGrid. fromEmail must be a verified sender in SendGrid.
// sendURL is the HTTP endpoint (e.g. from SENDGRID_API_URL); if empty, DefaultSendGridAPIURL is used.
func NewSendGridSender(apiKey, fromEmail, sendURL string) Sender {
	url := strings.TrimSpace(sendURL)
	if url == "" {
		url = DefaultSendGridAPIURL
	}
	return &SendGridSender{
		APIKey:    strings.TrimSpace(apiKey),
		FromEmail: strings.TrimSpace(fromEmail),
		SendURL:   url,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type sendGridMailBody struct {
	Personalizations []sendGridPersonalization `json:"personalizations"`
	From             sendGridAddress           `json:"from"`
	Subject          string                    `json:"subject"`
	Content          []sendGridContent         `json:"content"`
}

type sendGridPersonalization struct {
	To []sendGridAddress `json:"to"`
}

type sendGridAddress struct {
	Email string `json:"email"`
}

type sendGridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SendPasswordReset sends a plain-text message with the temporary password.
func (s *SendGridSender) SendPasswordReset(ctx context.Context, toEmail, plaintextTemporaryPassword string) error {
	if s.HTTPClient == nil {
		s.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}

	bodyText := fmt.Sprintf(`Olá,

Foi solicitada uma nova senha para a sua conta. Utilize a senha temporária abaixo para entrar; depois altere-a no site.

Senha temporária: %s

Se não foi você, ignore este e-mail.

— Equipe`, plaintextTemporaryPassword)

	payload := sendGridMailBody{
		Personalizations: []sendGridPersonalization{
			{To: []sendGridAddress{{Email: strings.TrimSpace(toEmail)}}},
		},
		From:    sendGridAddress{Email: s.FromEmail},
		Subject: "Redefinição de senha",
		Content: []sendGridContent{
			{Type: "text/plain", Value: bodyText},
		},
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sendgrid: encode body: %w", err)
	}

	url := strings.TrimSpace(s.SendURL)
	if url == "" {
		url = DefaultSendGridAPIURL
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("sendgrid: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("sendgrid: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("sendgrid: status %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
}
