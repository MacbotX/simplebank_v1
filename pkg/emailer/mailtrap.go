// Mailtrap implementation
package emailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	c "github.com/MacbotX/simplebank_v1/util"

)

type MailtrapMailer struct {
	ApiKey string
}

// Data structure passed into HTML template
type EmailTemplateData struct {
	Identifier string
	Token      string
}

func (m *MailtrapMailer) SendOTP(identifier string, token string) (string, error) {
	// build payload and send using Mailtrap
	subject := "[QicTik]-OTP Token"
	data := EmailTemplateData{
		Identifier: identifier,
		Token:      token,
	}

	// Generate HTML content
	html, err := GenerateHTML("./pkg/emailer/templates/otp.html", data)
	if err != nil {
		panic(err)
	}

	// Send OTP mailer
	err = Sender(identifier, subject, html)
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf(" OTP token sent to the identifier %s", identifier)
	return res, nil
}

// Send Welcome Message
func (m *MailtrapMailer) SendWelcomeMessage(identifier string) (string, error) {
	return "Mailtrap Welcome message sent", nil
}

// Email Sender
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type EmailPayload struct {
	From     EmailAddress   `json:"from"`
	To       []EmailAddress `json:"to"`
	Subject  string         `json:"subject"`
	HTML     string         `json:"html,omitempty"`
	Category string         `json:"category,omitempty"`
}

func Sender(identifier string, subject string, html string) error {
	// Load config data
	config, err := c.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	url := config.MailtrapURL
	method := "POST"

	// Prepare payload
	payload := EmailPayload{
		From: EmailAddress{
			Email: config.DefaultFromEmail,
			Name:  config.EmailSubjectPrefit,
		},
		To: []EmailAddress{
			{Email: identifier},
		},
		Subject: subject,
		HTML:    html,
	}
	// Marshal payload to json
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
		return err
	}
	// API Bearer token
	BearerToken := fmt.Sprintf("Bearer %s", config.MailtrapAuthToken)

	req.Header.Add("Authorization", BearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return err
}
