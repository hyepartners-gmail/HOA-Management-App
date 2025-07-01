package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type AgendaRequest struct {
	ID                string    `datastore:"id"`
	SubmittedByUserID string    `datastore:"submitted_by"`
	Subject           string    `datastore:"subject"`
	Description       string    `datastore:"description"`
	RequestedDate     time.Time `datastore:"requested_meeting_date"`
	CreatedAt         time.Time `datastore:"created_at"`
}

func SaveAgendaRequest(req AgendaRequest) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	req.CreatedAt = time.Now()

	key := datastore.NameKey("AgendaRequest", req.ID, nil)
	_, err := client.Put(ctx, key, &req)
	return err
}

func GetAllAgendaRequests() ([]*AgendaRequest, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var requests []*AgendaRequest
	query := datastore.NewQuery("AgendaRequest").Order("-created_at")
	_, err := client.GetAll(ctx, query, &requests)
	return requests, err
}
