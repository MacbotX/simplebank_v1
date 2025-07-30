// Decides which provider to use
package emailer

import (
	"fmt"

	c "github.com/MacbotX/simplebank_v1/util"
)

type Config struct {
	Provider          string
	MailtrapAuthToken string
	SendGridAuthToken string
}

func NewMailer(cfg c.Config) (Mailer, error) {
	switch cfg.Provider {
	case "mailtrap":
		return &MailtrapMailer{ApiKey: cfg.MailtrapAuthToken}, nil

	case "sendgrid":
		return &SendGridMailer{ApiKey: cfg.SendGridAuthToken}, nil

	default:
		return nil, fmt.Errorf("unsupported mailer provider: %s", cfg.Provider)
	}
}
