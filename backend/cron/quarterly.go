package cron

import (
	"log"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/google/uuid"
)

func RunQuarterlyBilling() error {
	now := time.Now()
	q1 := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	cabins, err := models.GetActiveCabins()
	if err != nil {
		return err
	}

	for _, cabin := range cabins {
		primaryOwner, err := models.GetPrimaryOwner(cabin.ID)
		if err != nil {
			log.Printf("Skipping cabin %s: %v", cabin.Label, err)
			continue
		}

		amountDue := float64(cabin.ShareCount) * 100 // example: $100 per share
		dueDate := q1.AddDate(0, 0, 30)

		invoice := models.Invoice{
			ID:              uuid.New(),
			CabinID:         cabin.ID,
			OwnerID:         primaryOwner.ID,
			PeriodStartDate: q1,
			PeriodEndDate:   q1.AddDate(0, 3, -1),
			AmountDue:       amountDue,
			DueDate:         dueDate,
			Status:          models.InvoiceSent,
			LateFeeApplied:  false,
			PaymentMethod:   models.PaymentNone,
		}

		pdfUrl, err := utils.GenerateInvoicePDF(invoice.ID)
		if err != nil {
			log.Printf("PDF generation failed for invoice %s: %v", invoice.ID, err)
			continue
		}
		invoice.PDFUrl = pdfUrl

		if err := models.SaveInvoice(invoice); err != nil {
			log.Printf("Failed to save invoice for cabin %s: %v", cabin.Label, err)
			continue
		}

		if err := utils.SendInvoiceEmail(primaryOwner.Email, invoice); err != nil {
			log.Printf("Failed to send invoice email to %s: %v", primaryOwner.Email, err)
		}
	}

	return nil
}
