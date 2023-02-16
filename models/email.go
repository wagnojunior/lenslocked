package models

import "github.com/go-mail/mail/v2"

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
