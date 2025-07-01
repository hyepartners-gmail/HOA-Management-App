package main

import (
	"time"

	"github.com/google/uuid"
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
}
