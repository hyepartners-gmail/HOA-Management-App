package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StripeCheckoutSession struct {
	URL       string
	SessionID string
	ExpiresAt time.Time
}

// Simulate a Stripe Checkout link
func CreateStripeCheckout(invoiceID uuid.UUID, amount float64, cabinLabel string) (*StripeCheckoutSession, error) {
	// Placeholder for real Stripe logic
	session := &StripeCheckoutSession{
		URL:       fmt.Sprintf("https://pay.stripe.com/session/%s", invoiceID.String()),
		SessionID: fmt.Sprintf("sess_%s", invoiceID.String()),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	return session, nil
}
