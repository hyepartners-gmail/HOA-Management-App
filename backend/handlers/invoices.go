package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// GetInvoicesHandler returns all invoices for admin/treasurer
func GetInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	invoices := models.GetAllInvoices()
	respondWithJSON(w, http.StatusOK, invoices)
}

// GetInvoiceByIDHandler fetches an invoice by ID
func GetInvoiceByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	invoice, err := models.GetInvoiceByID(id)
	if err != nil {
		utils.JSONError(w, "Invoice not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, invoice)
}

// CreateInvoiceHandler allows admin to manually create an invoice
func CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	var inv models.Invoice
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		utils.JSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	inv.ID = uuid.New().String()
	inv.Status = models.InvoiceStatusDraft
	inv.LateFeeApplied = false
	inv.CreatedAt = time.Now()
	models.SaveInvoice(inv)
	respondWithJSON(w, http.StatusCreated, inv)
}

// MarkInvoicePaidHandler updates invoice status and notes
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
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func GetInvoicePDFURLHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	invoiceID := chi.URLParam(r, "id")

	// Load invoice from Datastore
	invoice, err := models.GetInvoiceByID(invoiceID)
	if err != nil || invoice.OwnerID != user.AssociatedOwnerID {
		utils.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	url, err := utils.GenerateSignedInvoiceURL(invoice.PDFUrl)
	if err != nil {
		utils.JSONError(w, "failed to generate url", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
