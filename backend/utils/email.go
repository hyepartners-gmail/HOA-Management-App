// utils/email.go
package utils

import (
	"fmt"
	"os"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
)

func SendEmail(to, subject, body string) error {
	key := os.Getenv("SENDGRID_KEY")
	if key == "" {
		return fmt.Errorf("SENDGRID_KEY not set")
	}
	// Stubbed logic
	fmt.Printf("Sending email to %s with subject: %s\n%s\n", to, subject, body)
	return nil
}

func SendNotificationEmail(n models.Notification) {
	recipients := models.ResolveRecipients(n)
	for _, email := range recipients.Emails {
		go SendEmail(
			email,
			fmt.Sprintf("HOA Notification: %s", n.Title),
			n.Body,
		)
	}
}

func SendNewsletterToAllOwners(subject, body string) {
	owners, _ := models.GetAllOwners() // You already have this
	for _, o := range owners {
		if o.Email != "" {
			_ = SendEmail(o.Email, subject, body)
		}
	}
}
