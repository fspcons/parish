package email

import (
	"context"
	"log/slog"
)

// LogSender logs password-reset e-mails (development / tests). Do not use in production.
type LogSender struct{}

// SendPasswordReset writes the intent to structured logs; the plaintext password is included for local debugging only.
func (LogSender) SendPasswordReset(_ context.Context, toEmail, plaintextTemporaryPassword string) error {
	slog.Info("password reset email (log sender)",
		"to", toEmail,
		"tempPasswordLength", len(plaintextTemporaryPassword),
	)
	return nil
}
