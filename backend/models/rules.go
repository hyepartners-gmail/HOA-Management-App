package models

import (
	"context"

	"cloud.google.com/go/datastore"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type HOARules struct {
	ID            string `datastore:"id"` // single row w/ known ID like "default"
	Renovation    string `datastore:"renovation"`
	PaintColors   string `datastore:"paint_colors"`
	ShingleTypes  string `datastore:"shingle_types"`
	GeneralBylaws string `datastore:"general_bylaws"`
	LastUpdatedBy string `datastore:"last_updated_by"`
}

func GetHOARules() (*HOARules, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var rules HOARules
	key := datastore.NameKey("HOARules", "default", nil)
	if err := client.Get(ctx, key, &rules); err != nil {
		return nil, err
	}
	return &rules, nil
}

func SaveHOARules(rules HOARules) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("HOARules", "default", nil)
	_, err := client.Put(ctx, key, &rules)
	return err
}
