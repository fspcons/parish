package email

import "context"

// Sender delivers transactional e-mail (e.g. password reset).
type Sender interface {
	SendPasswordReset(ctx context.Context, toEmail, plaintextTemporaryPassword string) error
}
