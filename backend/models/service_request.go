package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type ServiceRequest struct {
	ID                uuid.UUID `datastore:"id"`
	Category          string    `datastore:"category"`
	SubmittedByUserID uuid.UUID `datastore:"submitted_by_user_id"`
	Description       string    `datastore:"description"`
	CreatedAt         time.Time `datastore:"created_at"`
	Status            string    `datastore:"status"` // open, resolved, escalated
}

func SaveServiceRequest(req ServiceRequest) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)
	key := datastore.NameKey("ServiceRequest", req.ID.String(), nil)
	_, err := client.Put(ctx, key, &req)
	return err
}

func GetAllServiceRequests() ([]*ServiceRequest, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var results []*ServiceRequest
	query := datastore.NewQuery("ServiceRequest").Order("-created_at")
	_, err := client.GetAll(ctx, query, &results)
	return results, err
}

func UpdateServiceRequestStatus(id string, status string) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("ServiceRequest", id, nil)
	var req ServiceRequest
	if err := client.Get(ctx, key, &req); err != nil {
		return err
	}
	req.Status = status
	_, err := client.Put(ctx, key, &req)
	return err
}
