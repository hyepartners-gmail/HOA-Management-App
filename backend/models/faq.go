package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type FAQ struct {
	ID        uuid.UUID `datastore:"id" json:"id"`
	Title     string    `datastore:"title" json:"title"`
	Content   string    `datastore:"content" json:"content"`
	CreatedAt time.Time `datastore:"created_at" json:"created_at"`
	UpdatedAt time.Time `datastore:"updated_at" json:"updated_at"`
}

func SaveFAQ(f FAQ) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)
	key := datastore.NameKey("FAQ", f.ID.String(), nil)
	_, err := client.Put(ctx, key, &f)
	return err
}

func GetAllFAQs() ([]*FAQ, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)
	var faqs []*FAQ
	query := datastore.NewQuery("FAQ").Order("title")
	_, err := client.GetAll(ctx, query, &faqs)
	return faqs, err
}
