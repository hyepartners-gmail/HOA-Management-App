package models

import (
	"context"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
)

type Document struct {
	ID         uuid.UUID `datastore:"id" json:"id"`
	Title      string    `datastore:"title" json:"title"`
	Category   string    `datastore:"category" json:"category"`
	URL        string    `datastore:"url" json:"url"`
	VisibleTo  string    `datastore:"visible_to" json:"visible_to"` // e.g. "admin", "board", "owners", "public"
	UploadedBy uuid.UUID `datastore:"uploaded_by" json:"uploaded_by"`
	UploadedAt time.Time `datastore:"uploaded_at" json:"uploaded_at"`
}

func SaveDocument(doc Document) error {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)

	key := datastore.NameKey("Document", doc.ID.String(), nil)
	_, err := client.Put(ctx, key, &doc)
	return err
}

func ListDocuments(role string) ([]Document, error) {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)

	q := datastore.NewQuery("Document").Order("-UploadedAt").Filter("VisibleTo >=", role)

	var results []Document
	_, err := client.GetAll(ctx, q, &results)
	return results, err
}
