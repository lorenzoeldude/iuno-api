package email

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

func SendVerificationEmail(to string, token string) error {
	verificationURL := fmt.Sprintf(
		"https://www.iunoni.com/verify-email?token=%s",
		token,
	)

	params := &resend.SendEmailRequest{
		From:    "Iunoni <noreply@iunoni.com>",
		To:      []string{to},
		Subject: "Verify your email",
		Html: fmt.Sprintf(`
			<h2>Welcome to Iunoni!</h2>

			<p>Thanks for signing up!</p>

			<p>Please verify your email by clicking the button below.</p>

			<p>
				<a href="%s">Verify Email</a>
			</p>

			<p>This link expires in 24 hours.</p>
		`, verificationURL),
	}

	_, err := Client.Emails.Send(params)
	return err
}