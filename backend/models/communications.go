package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type CommunicationType string

const (
	CommNewsletter   CommunicationType = "newsletter"
	CommAgenda       CommunicationType = "agenda"
	CommMinutes      CommunicationType = "minutes"
	CommFlash        CommunicationType = "flash"
	CommAnnouncement CommunicationType = "announcement"
)

type Communication struct {
	ID        uuid.UUID         `datastore:"id" json:"id"`
	Title     string            `datastore:"title" json:"title"`
	Body      string            `datastore:"body,noindex" json:"body"`
	Type      CommunicationType `datastore:"type" json:"type"`
	CreatedBy uuid.UUID         `datastore:"created_by" json:"created_by"`
	CreatedAt time.Time         `datastore:"created_at" json:"created_at"`
}

func SaveCommunication(c Communication) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("Communication", c.ID.String(), nil)
	_, err := client.Put(ctx, key, &c)
	return err
}

func ListCommunications(limit int, commType string) ([]Communication, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	q := datastore.NewQuery("Communication").Order("-CreatedAt").Limit(limit)
	if commType != "" {
		q = q.Filter("Type =", commType)
	}

	var results []Communication
	_, err := client.GetAll(ctx, q, &results)
	return results, err
}

func GetCommunicationByID(id string) (*Communication, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var c Communication
	key := datastore.NameKey("Communication", id, nil)
	err := client.Get(ctx, key, &c)
	return &c, err
}
