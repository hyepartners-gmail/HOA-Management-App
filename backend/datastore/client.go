package datastore

import (
	"context"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/datastore"
)

var (
	client     *datastore.Client
	clientOnce sync.Once
)

func GetClient(ctx context.Context) *datastore.Client {
	clientOnce.Do(func() {
		var err error
		client, err = datastore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT")) // ← Replace with env var
		if err != nil {
			log.Fatalf("Failed to create datastore client: %v", err)
		}
	})
	return client
}
