// models/store.go
package models

import (
	"context"

	"cloud.google.com/go/datastore"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

// Save any entity to Datastore
func saveToDatastore(entity interface{}, kind string, id string) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)
	key := datastore.NameKey(kind, id, nil)
	_, err := client.Put(ctx, key, entity)
	return err
}

// Load entity by ID
func loadFromDatastore(dst interface{}, kind string, id string) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)
	key := datastore.NameKey(kind, id, nil)
	return client.Get(ctx, key, dst)
}

// Delete entity by ID
func deleteFromDatastore(kind string, id string) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)
	key := datastore.NameKey(kind, id, nil)
	return client.Delete(ctx, key)
}

// Delete entities by field value (e.g., delete all comments by post_id)
func deleteByField(kind string, field string, value string) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	query := datastore.NewQuery(kind).Filter(field+" =", value).KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		return err
	}

	return client.DeleteMulti(ctx, keys)
}
