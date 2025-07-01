package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateInvoicePDF(invoice models.Invoice, cabin models.Cabin, owner models.Owner) (string, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.Cell(40, 10, "Bear Paw Cabins HOA")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice ID: %s", invoice.ID.String()))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", invoice.PeriodStartDate.Format("Jan 2, 2006")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("Jan 2, 2006")))
	pdf.Ln(12)

	// Owner Info
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Billed To:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, owner.FullName)
	pdf.Ln(6)
	pdf.MultiCell(0, 6, owner.MailingAddress, "", "", false)
	pdf.Ln(10)

	// Cabin Info
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Cabin Details:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Cabin: %s", cabin.Label))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Shares: %d", cabin.ShareCount))
	pdf.Ln(10)

	// Amount Due
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Invoice Summary:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Amount Due: $%.2f", invoice.AmountDue))
	if invoice.LateFeeApplied {
		pdf.Ln(6)
		pdf.Cell(40, 10, "(Late fee applied)")
	}
	pdf.Ln(10)

	// Notes
	if invoice.Notes != "" {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Notes:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(0, 6, invoice.Notes, "", "", false)
	}

	// Save locally (e.g., /tmp/invoices/<invoice-id>.pdf)
	filename := fmt.Sprintf("%s.pdf", invoice.ID.String())
	outDir := "/tmp/invoices"
	os.MkdirAll(outDir, 0755)
	fullPath := filepath.Join(outDir, filename)

	err := pdf.OutputFileAndClose(fullPath)
	if err != nil {
		return "", err
	}

	gcsPath, err := UploadFileToGCS(fullPath, filename)
	if err != nil {
		return "", err
	}
	return gcsPath, nil
}
