package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type MeetingMinutes struct {
	ID               string    `datastore:"id"`
	MeetingDate      time.Time `datastore:"meeting_date"`
	Title            string    `datastore:"title"`
	ContentHTML      string    `datastore:"content_html"` // or URL to uploaded PDF
	Published        bool      `datastore:"published"`
	UploadedByUserID string    `datastore:"uploaded_by_user_id"`
	CreatedAt        time.Time `datastore:"created_at"`
}

func SaveMeetingMinutes(m MeetingMinutes) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	m.CreatedAt = time.Now()

	key := datastore.NameKey("MeetingMinutes", m.ID, nil)
	_, err := client.Put(ctx, key, &m)
	return err
}

func GetAllMeetingMinutes() ([]*MeetingMinutes, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var minutes []*MeetingMinutes
	_, err := client.GetAll(ctx, datastore.NewQuery("MeetingMinutes").Order("-MeetingDate"), &minutes)
	return minutes, err
}
