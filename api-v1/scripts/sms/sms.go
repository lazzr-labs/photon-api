package sms

import (
	"errors"
	"strings"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

var twilioAccountSID string
var twilioAuthToken string
var twilioFromPhone string

func SetClient(accountSID string, authToken string, fromPhone string) {
	twilioAccountSID = accountSID
	twilioAuthToken = authToken
	twilioFromPhone = fromPhone
}

func SendMessage(toPhone string, body string) error {
	if strings.TrimSpace(twilioAccountSID) == "" || strings.TrimSpace(twilioAuthToken) == "" || strings.TrimSpace(twilioFromPhone) == "" {
		return errors.New("twilio sms client is not configured")
	}
	if strings.TrimSpace(toPhone) == "" {
		return errors.New("sms recipient phone is required")
	}
	if strings.TrimSpace(body) == "" {
		return errors.New("sms message body is required")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioAccountSID,
		Password: twilioAuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(twilioFromPhone)
	params.SetTo(toPhone)
	params.SetBody(body)

	_, err := client.Api.CreateMessage(params)
	return err
}
