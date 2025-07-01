package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

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

// SendSMSTwilio sends an SMS using the Twilio API
func SendSMSTwilio(to string, body string) error {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")

	if accountSID == "" || authToken == "" || fromNumber == "" {
		return fmt.Errorf("Twilio credentials are not fully set in environment variables")
	}

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", fromNumber)
	msgData.Set("Body", body)
	reqBody := *strings.NewReader(msgData.Encode())

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)

	req, _ := http.NewRequest("POST", urlStr, &reqBody)
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Twilio API returned status %s", resp.Status)
	}

	return nil
}
