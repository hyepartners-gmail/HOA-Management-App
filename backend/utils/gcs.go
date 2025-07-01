package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// UploadFileToGCS uploads a local file to GCS and returns the signed URL.
func UploadFileToGCS(localPath string, filename string) (string, error) {
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

	bucket := client.Bucket(GCSBucketName)
	objectPath := filepath.Join(InvoiceFolderPath, filename)
	writer := bucket.Object(objectPath).NewWriter(ctx)
	writer.ContentType = "application/pdf"

	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to write to GCS: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Optionally return a signed URL
	return GenerateSignedInvoiceURL(filename)
}

// UploadToGCS uploads bytes to GCS and returns the public or signed URL.
func UploadToGCS(data []byte, objectName string) (string, error) {
	ctx := context.Background()

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

	bucket := client.Bucket(GCSBucketName)
	writer := bucket.Object(objectName).NewWriter(ctx)
	writer.ContentType = "application/octet-stream"

	if _, err := writer.Write(data); err != nil {
		return "", fmt.Errorf("failed to write to GCS: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	return GenerateSignedInvoiceURL(objectName)
}
