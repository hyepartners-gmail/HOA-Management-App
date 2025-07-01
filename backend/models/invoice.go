package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

// InvoiceStatus represents the current state of the invoice.
type InvoiceStatus string

const (
	InvoiceDraft InvoiceStatus = "draft"
	InvoiceSent  InvoiceStatus = "sent"
	InvoicePaid  InvoiceStatus = "paid"
	InvoiceLate  InvoiceStatus = "late"
)

// PaymentMethod represents how the invoice was paid.
type PaymentMethod string

const (
	PaymentStripe PaymentMethod = "stripe"
	PaymentCheck  PaymentMethod = "check"
	PaymentZelle  PaymentMethod = "zelle"
	PaymentNone   PaymentMethod = "none"
)

type Invoice struct {
	ID              uuid.UUID     `datastore:"id" json:"id"`
	CabinID         uuid.UUID     `datastore:"cabin_id" json:"cabin_id"`
	OwnerID         uuid.UUID     `datastore:"owner_id" json:"owner_id"`
	PeriodStartDate time.Time     `datastore:"period_start_date" json:"period_start_date"`
	PeriodEndDate   time.Time     `datastore:"period_end_date" json:"period_end_date"`
	AmountDue       float64       `datastore:"amount_due" json:"amount_due"`
	DueDate         time.Time     `datastore:"due_date" json:"due_date"`
	Status          InvoiceStatus `datastore:"status" json:"status"`
	LateFeeApplied  bool          `datastore:"late_fee_applied" json:"late_fee_applied"`
	PDFUrl          string        `datastore:"pdf_url" json:"pdf_url"`
	PaymentMethod   PaymentMethod `datastore:"payment_method" json:"payment_method"`
	PaidAt          *time.Time    `datastore:"paid_at" json:"paid_at,omitempty"`
	Notes           string        `datastore:"notes" json:"notes,omitempty"`
	CreatedAt       time.Time     `datastore:"created_at" json:"created_at"` // Add this
}

func GetAllInvoices() []*Invoice {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var invoices []*Invoice
	_, _ = client.GetAll(ctx, datastore.NewQuery("Invoice").Order("-due_date"), &invoices)
	return invoices
}

func GetInvoiceByID(id string) (*Invoice, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key, err := datastore.DecodeKey(id)
	if err != nil {
		return nil, err
	}

	var invoice Invoice
	if err := client.Get(ctx, key, &invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

func SaveInvoice(inv Invoice) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("Invoice", inv.ID.String(), nil)
	_, err := client.Put(ctx, key, &inv)
	return err
}

type ManualPaymentUpdate struct {
	PaymentDate time.Time
	Method      PaymentMethod
	Notes       string
}

func UpdateInvoicePayment(invoiceID string, update ManualPaymentUpdate) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("Invoice", invoiceID, nil)
	var invoice Invoice
	if err := client.Get(ctx, key, &invoice); err != nil {
		return err
	}

	invoice.PaymentMethod = update.Method
	invoice.PaidAt = &update.PaymentDate
	invoice.Notes = update.Notes
	invoice.Status = InvoicePaid

	_, err := client.Put(ctx, key, &invoice)
	return err
}

func GenerateQuarterlyInvoices() error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	// Placeholder logic: generate one invoice per active cabin
	var cabins []*Cabin
	_, err := client.GetAll(ctx, datastore.NewQuery("Cabin").Filter("is_active =", true), &cabins)
	if err != nil {
		return err
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month()-3, 1, 0, 0, 0, 0, time.UTC)
	end := now

	for _, cabin := range cabins {
		inv := Invoice{
			ID:              uuid.New(),
			CabinID:         uuid.MustParse(cabin.ID),
			OwnerID:         uuid.MustParse(cabin.PrimaryOwnerID),
			PeriodStartDate: start,
			PeriodEndDate:   end,
			AmountDue:       500.00, // example static rate
			DueDate:         end.AddDate(0, 0, 30),
			Status:          InvoiceDraft,
			CreatedAt:       time.Now(),
		}
		key := datastore.NameKey("Invoice", inv.ID.String(), nil)
		if _, err := client.Put(ctx, key, &inv); err != nil {
			return err
		}
	}

	return nil
}
