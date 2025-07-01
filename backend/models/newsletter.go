package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type Newsletter struct {
	ID              uuid.UUID  `datastore:"id" json:"id"`
	Title           string     `datastore:"title" json:"title"`
	Body            string     `datastore:"body,noindex" json:"body"` // long-form HTML/text
	PublishedAt     *time.Time `datastore:"published_at" json:"published_at,omitempty"`
	CreatedByUserID uuid.UUID  `datastore:"created_by_user_id" json:"created_by_user_id"`
}

func SaveNewsletter(n *Newsletter) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}

	key := datastore.NameKey("Newsletter", n.ID.String(), nil)
	_, err := client.Put(ctx, key, n)
	return err
}

func GetAllNewsletters() ([]Newsletter, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var newsletters []Newsletter
	query := datastore.NewQuery("Newsletter").Order("-published_at")
	_, err := client.GetAll(ctx, query, &newsletters)
	return newsletters, err
}

func GetNewsletterByID(id string) (*Newsletter, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var n Newsletter
	key := datastore.NameKey("Newsletter", id, nil)
	if err := client.Get(ctx, key, &n); err != nil {
		return nil, err
	}
	return &n, nil
}
