// handlers/invoice.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/google/uuid"
)

// GET /api/invoices
func GetInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	// role: treasurer or cabin_owner
	// retrieve and return invoices based on role/ownership
}

// POST /api/invoices/manual-payment
func RecordManualPaymentHandler(w http.ResponseWriter, r *http.Request) {
	type PaymentInput struct {
		InvoiceID   uuid.UUID            `json:"invoice_id"`
		PaymentDate time.Time            `json:"payment_date"`
		Method      models.PaymentMethod `json:"method"`
		Notes       string               `json:"notes"`
	}

	var input PaymentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "invalid input", http.StatusBadRequest)
		return
	}

	// lookup invoice, validate, apply update, persist
	// set status to paid and update fields accordingly
	w.WriteHeader(http.StatusNoContent)
}

// POST /api/invoices/trigger-quarterly
func GenerateQuarterlyInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	// triggered by cron job or admin button
	// iterate all active cabins
	// calculate amount based on shares
	// create invoice
	// attach pdf (optional stubbed)
	// queue for email delivery
	w.WriteHeader(http.StatusNoContent)
}
