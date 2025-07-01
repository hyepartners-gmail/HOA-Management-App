package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type AuditLog struct {
	ID                uuid.UUID `datastore:"id" json:"id"`
	ActionType        string    `datastore:"action_type" json:"action_type"`
	PerformedByUserID uuid.UUID `datastore:"performed_by_user_id" json:"performed_by_user_id"`
	TargetID          string    `datastore:"target_id" json:"target_id"`
	TargetType        string    `datastore:"target_type" json:"target_type"`
	Timestamp         time.Time `datastore:"timestamp" json:"timestamp"`
	Meta              string    `datastore:"meta,noindex" json:"meta"`
}

func LogAction(log AuditLog) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	log.ID = uuid.New()
	log.Timestamp = time.Now()
	key := datastore.NameKey("AuditLog", log.ID.String(), nil)
	_, err := client.Put(ctx, key, &log)
	return err
}

func ListAuditLogs(limit int) ([]AuditLog, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	q := datastore.NewQuery("AuditLog").Order("-Timestamp").Limit(limit)
	var logs []AuditLog
	_, err := client.GetAll(ctx, q, &logs)
	return logs, err
}
