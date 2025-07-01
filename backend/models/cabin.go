// models/cabin.go
package models

import (
	"context"

	"cloud.google.com/go/datastore"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type Cabin struct {
	ID             string   `datastore:"id"`
	Label          string   `datastore:"label"`
	BedroomCount   int      `datastore:"bedroom_count"`
	ShareCount     int      `datastore:"share_count"`
	OwnerIDs       []string `datastore:"owners"`
	PrimaryOwnerID string   `datastore:"primary_owner_id"`
	IsActive       bool     `datastore:"is_active"`
}

func GetAllCabins() ([]*Cabin, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var cabins []*Cabin
	query := datastore.NewQuery("Cabin")
	_, err := client.GetAll(ctx, query, &cabins)
	return cabins, err
}

func SaveCabin(c *Cabin) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)
	key := datastore.NameKey("Cabin", c.ID, nil)
	_, err := client.Put(ctx, key, c)
	return err
}

func UpdateCabin(c *Cabin) error {
	return SaveCabin(c)
}

func FindCabinByID(id string) (*Cabin, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var c Cabin
	key := datastore.NameKey("Cabin", id, nil)
	if err := client.Get(ctx, key, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
