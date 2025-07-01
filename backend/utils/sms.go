package utils

import (
	"fmt"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
)

func SendFlashSMS(n models.Notification) {
	shortBody := n.Body
	if len(shortBody) > 300 {
		shortBody = shortBody[:297] + "..."
	}
	recipients := models.ResolveRecipients(n)
	for _, phone := range recipients.Phones {
		go SendSMSTwilio(phone, fmt.Sprintf("%s: %s", n.Title, shortBody))
	}
}
