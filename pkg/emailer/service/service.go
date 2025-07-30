package emailer

import (
	"log"

	c "github.com/MacbotX/simplebank_v1/util"

	"github.com/MacbotX/simplebank_v1/pkg/emailer"
)

func SendEmailOTP(identifier string, token string) (string, error) {
	// central config (ideally loaded from .env)
	config, err := c.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	cfg := c.Config{
		Provider:          config.Provider,
		MailtrapAuthToken: config.MailtrapAuthToken,
		SendGridAuthToken: config.SendGridAuthToken,
	}

	mailer, err := emailer.NewMailer(cfg)
	if err != nil {
		return "", err
	}
	res, err := mailer.SendOTP(identifier, token)
	return res, err
}
