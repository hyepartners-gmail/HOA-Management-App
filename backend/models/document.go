package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type Document struct {
	ID         uuid.UUID `datastore:"id" json:"id"`
	Title      string    `datastore:"title" json:"title"`
	Category   string    `datastore:"category" json:"category"`
	URL        string    `datastore:"url" json:"url"`
	VisibleTo  string    `datastore:"visible_to" json:"visible_to"`
	UploadedBy string    `datastore:"uploaded_by" json:"uploaded_by"` // <- updated
	UploadedAt time.Time `datastore:"uploaded_at" json:"uploaded_at"`
}

func SaveDocument(doc Document) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("Document", doc.ID.String(), nil)
	_, err := client.Put(ctx, key, &doc)
	return err
}

func ListDocuments(role string) ([]Document, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	q := datastore.NewQuery("Document").Order("-UploadedAt").Filter("VisibleTo >=", role)

	var results []Document
	_, err := client.GetAll(ctx, q, &results)
	return results, err
}
