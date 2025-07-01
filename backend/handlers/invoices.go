package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

// GET /api/invoices - treasurer or cabin_owner
func GetInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	invoices := models.GetAllInvoices() // add filtering by role/user later if needed
	utils.RespondWithJSON(w, http.StatusOK, invoices)
}

// GET /api/invoices/{id}
func GetInvoiceByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	invoice, err := models.GetInvoiceByID(id)
	if err != nil {
		utils.JSONError(w, "Invoice not found", http.StatusNotFound)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, invoice)
}

// POST /api/invoices/manual-payment
func RecordManualPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		InvoiceID   uuid.UUID            `json:"invoice_id"`
		PaymentDate time.Time            `json:"payment_date"`
		Method      models.PaymentMethod `json:"method"`
		Notes       string               `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "invalid input", http.StatusBadRequest)
		return
	}
	err := models.UpdateInvoicePayment(input.InvoiceID.String(), models.ManualPaymentUpdate{
		PaymentDate: input.PaymentDate,
		Method:      input.Method,
		Notes:       input.Notes,
	})
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /api/invoices
func CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	var inv models.Invoice
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		utils.JSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	inv.ID = uuid.New()
	inv.Status = models.InvoiceDraft
	inv.CreatedAt = time.Now()
	err := models.SaveInvoice(inv)
	if err != nil {
		utils.JSONError(w, "failed to save invoice", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, inv)
}

// POST /api/invoices/{id}/mark-paid
func MarkInvoicePaidHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var update models.ManualPaymentUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err := models.UpdateInvoicePayment(id, update)
	if err != nil {
		utils.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// GET /api/invoices/{id}/pdf-url
func GetInvoicePDFURLHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	invoiceID := chi.URLParam(r, "id")

	invoice, err := models.GetInvoiceByID(invoiceID)
	if err != nil || invoice.OwnerID.String() != user.AssociatedOwnerID {
		utils.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	url, err := utils.GenerateSignedInvoiceURL(invoice.PDFUrl)
	if err != nil {
		utils.JSONError(w, "failed to generate url", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"url": url})
}

func GenerateQuarterlyInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	err := models.GenerateQuarterlyInvoices()
	if err != nil {
		utils.JSONError(w, "Failed to generate invoices", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
