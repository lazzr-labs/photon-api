package emails

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var sendgridApiKey string

var fromName = "Lazzr"
var fromEmail = "support@lazzr.com"

func SetApiKey(ApiKey string) {
	sendgridApiKey = ApiKey
}

func SendEmail(email *mail.SGMailV3) {
	request := sendgrid.GetRequest(sendgridApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Body = mail.GetRequestBody(email)
	request.Method = "POST"
	sendgrid.API(request)
}
