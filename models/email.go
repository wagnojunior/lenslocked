package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

// SMTPConfig defines a new type to hold the SMTP configuration. These credentials are provided by
// the SMTP provider (i.e.: mailtrap)
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// EmailService defines a new type to hold the EmailService configuration, specially the dialer.
type EmailService struct {
	// DefaultSender is used as the default sender when one isn't provided for an email. This is
	// also used in functions where the email is a predetermined, like the forgotten password email
	DefaultSender string

	// Unexported fields
	dialer *mail.Dialer
}

// NewEmailService constructs a new email service with the provided SMTP configuration
func NewEmailService(config SMTPConfig) *EmailService {
	sd := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}

	return &sd
}

// Email defines a new type that holdes the relevant information of an email. This type wraps the
// `mail.Message` type from the `mail` package so that users of this service do not have to worry
// about third-party dependency
type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}

// Send sends an single email and returns an error, if any
func (es *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	es.setFrom(msg, email)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)

	// Sets the email bode depending on whether a plain text or HTML are passed. In case either is
	// passed there is no need to add an alternative
	switch {
	case email.PlainText != "" && email.HTML != "":
		msg.SetBody("text/plain", email.PlainText)
		msg.AddAlternative("text/html", email.HTML)
	case email.PlainText != "":
		msg.SetBody("text/plain", email.PlainText)
	case email.HTML != "":
		msg.SetBody("text/hmtl", email.HTML)
	}

	// Sends a single email
	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

// ForgotPassword sends a reset-password email to a user
func (es *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		To:        to,
		Subject:   "Reset your password",
		PlainText: "To reset your password, please visit the following link:" + resetURL,
		HTML:      `<p>To reset your password, please visit the following link: <a href="` + resetURL + `">` + resetURL + `</a></p>`,
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}

	return nil
}

// setFrom sets the `from` field in an email.
func (es *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	// In case the FROM field is specified, use it
	case email.From != "":
		from = email.From
	// In case the FROM field is not specified and the email service's default sendere it, use it
	case es.DefaultSender != "":
		from = es.DefaultSender
	// In case the FROM field is not specified, use the constant default sender
	default:
		from = DefaultSender
	}

	msg.SetHeader("From", from)
}
