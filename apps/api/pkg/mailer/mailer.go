// Package mailer sends transactional email on behalf of usecases (e.g.
// AuthUsecase.ForgotPassword). It exposes EmailSender so usecases depend on
// a small interface rather than a specific provider's SDK, per
// docs/adr/0002-clean-architecture-backend.md.
package mailer

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

// EmailSender sends a single plain-text/HTML email. body is treated as
// HTML — callers that need a plain-text body should set one accordingly.
type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

// ResendSender implements EmailSender on top of Resend
// (https://resend.com).
type ResendSender struct {
	client *resend.Client
	from   string
}

// NewResendSender constructs a ResendSender. apiKey is ops-supplied (see
// RESEND_API_KEY in .env.example) with no default, same as the JWT
// secrets. from is the address emails are sent as (see
// EMAIL_FROM_ADDRESS).
func NewResendSender(apiKey, from string) *ResendSender {
	return &ResendSender{client: resend.NewClient(apiKey), from: from}
}

// Send sends body as the HTML content of an email to to.
func (s *ResendSender) Send(ctx context.Context, to, subject, body string) error {
	_, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	})
	if err != nil {
		return fmt.Errorf("sending email via resend: %w", err)
	}
	return nil
}
