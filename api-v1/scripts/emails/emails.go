package emails

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var ErrClientNotConfigured = errors.New("email client is not configured")
var ErrFromRequired = errors.New("email from name and email are required")
var ErrRecipientRequired = errors.New("email recipient is required")
var ErrTemplateRequired = errors.New("email template ID is required")

var sendgridApiKey string
var fromName string
var fromEmail string

type TemplateEmail struct {
	ToEmail     string
	ToName      string
	TemplateID  string
	Data        map[string]any
	Attachments []Attachment
}

type Attachment struct {
	Content     []byte
	Filename    string
	ContentType string
}

func SetClient(apiKey string, fromName string, fromEmail string) error {
	SetApiKey(apiKey)
	SetFrom(fromName, fromEmail)

	if strings.TrimSpace(sendgridApiKey) == "" {
		return nil
	}
	if strings.TrimSpace(fromName) == "" || strings.TrimSpace(fromEmail) == "" {
		return ErrFromRequired
	}

	return nil
}

func SetApiKey(apiKey string) {
	sendgridApiKey = strings.TrimSpace(apiKey)
}

func SetFrom(name string, email string) {
	fromName = strings.TrimSpace(name)
	fromEmail = strings.TrimSpace(email)
}

func SendTemplate(input TemplateEmail) error {
	if strings.TrimSpace(sendgridApiKey) == "" {
		return ErrClientNotConfigured
	}
	if strings.TrimSpace(fromName) == "" || strings.TrimSpace(fromEmail) == "" {
		return ErrFromRequired
	}

	input.ToEmail = strings.TrimSpace(input.ToEmail)
	input.TemplateID = strings.TrimSpace(input.TemplateID)
	if input.ToEmail == "" {
		return ErrRecipientRequired
	}
	if input.TemplateID == "" {
		return ErrTemplateRequired
	}

	email := mail.NewV3Mail()
	email.SetFrom(mail.NewEmail(fromName, fromEmail))
	email.SetTemplateID(input.TemplateID)

	personalization := mail.NewPersonalization()
	personalization.AddTos(mail.NewEmail(input.ToName, input.ToEmail))
	for key, value := range input.Data {
		personalization.SetDynamicTemplateData(key, value)
	}
	email.AddPersonalizations(personalization)

	for _, attachment := range input.Attachments {
		if len(attachment.Content) == 0 || strings.TrimSpace(attachment.Filename) == "" {
			continue
		}

		email.AddAttachment(attachment.ToSendGrid())
	}

	return SendEmail(email)
}

func SendEmail(email *mail.SGMailV3) error {
	if strings.TrimSpace(sendgridApiKey) == "" {
		return ErrClientNotConfigured
	}
	if email == nil {
		return errors.New("email message is required")
	}

	request := sendgrid.GetRequest(sendgridApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Body = mail.GetRequestBody(email)
	request.Method = "POST"

	response, err := sendgrid.API(request)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("sendgrid send failed: status=%d body=%s", response.StatusCode, response.Body)
	}

	return nil
}

func (attachment Attachment) ToSendGrid() *mail.Attachment {
	contentType := strings.TrimSpace(attachment.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	emailAttachment := mail.NewAttachment()
	emailAttachment.SetContent(base64.StdEncoding.EncodeToString(attachment.Content))
	emailAttachment.SetType(contentType)
	emailAttachment.SetFilename(attachment.Filename)
	emailAttachment.SetDisposition("attachment")
	return emailAttachment
}
