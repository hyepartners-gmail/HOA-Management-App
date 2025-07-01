package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	GCSBucketName        = "bearpaw-invoices"
	InvoiceFolderPath    = "invoices"
	GCSSignedURLDuration = 24 * time.Hour
)

// GenerateSignedInvoiceURL generates a secure signed URL for a given PDF object.
func GenerateSignedInvoiceURL(objectName string) (string, error) {
	ctx := context.Background()

	// Load credentials from environment
	serviceAccountKeyPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if serviceAccountKeyPath == "" {
		return "", fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	opts := option.WithCredentialsFile(serviceAccountKeyPath)
	client, err := storage.NewClient(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	url, err := storage.SignedURL(GCSBucketName, fmt.Sprintf("%s/%s", InvoiceFolderPath, objectName), &storage.SignedURLOptions{
		Method:         "GET",
		Expires:        time.Now().Add(GCSSignedURLDuration),
		ContentType:    "application/pdf",
		GoogleAccessID: os.Getenv("GCS_SIGNING_EMAIL"),       // from service account
		PrivateKey:     []byte(os.Getenv("GCS_PRIVATE_KEY")), // base64 decoded if needed
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return url, nil
}
