package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type ProxyAssignment struct {
	ID          string     `datastore:"id"`
	FromUserID  string     `datastore:"from_user_id"`
	ToUserID    *string    `datastore:"to_user_id,omitempty"` // optional if to_office used
	ToOffice    *string    `datastore:"to_office,omitempty"`  // enum: president, secretary, etc
	MeetingDate *time.Time `datastore:"meeting_date,omitempty"`
	IsOneTime   bool       `datastore:"is_one_time"`
	CreatedAt   time.Time  `datastore:"created_at"`
}

func SaveProxy(p ProxyAssignment) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	p.CreatedAt = time.Now()

	key := datastore.NameKey("ProxyAssignment", p.ID, nil)
	_, err := client.Put(ctx, key, &p)
	return err
}

func GetAllProxies() ([]ProxyAssignment, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var proxies []ProxyAssignment
	query := datastore.NewQuery("ProxyAssignment").Order("-created_at")
	_, err := client.GetAll(ctx, query, &proxies)
	return proxies, err
}
