package emails

import "github.com/sendgrid/sendgrid-go/helpers/mail"

var forgotPasswordTemplateID = "d-557450e3966b48fea0c91cf1075a01f3"

func ForgotPasswordEmail(userEmail string, userFirstName string, code string) {
	from := mail.NewEmail(fromName, fromEmail)
	tos := []*mail.Email{
		mail.NewEmail("", userEmail),
	}

	email := mail.NewV3Mail()
	email.SetFrom(from)
	email.SetTemplateID(forgotPasswordTemplateID)

	personalization := mail.NewPersonalization()
	personalization.AddTos(tos...)

	personalization.SetDynamicTemplateData("first_name", userFirstName)
	personalization.SetDynamicTemplateData("code", code)

	email.AddPersonalizations(personalization)
	SendEmail(email)
}
