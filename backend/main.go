// cmd/server/main.go  (or replace the existing main.go)
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/seed" // <— new
)

func main() {
	/* ---------------------------------------------------------
	   1. One-time seeding (runs only if each Kind is empty)
	--------------------------------------------------------- */
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT")         // Cloud Run sets this automatically
	namespace := os.Getenv("DATASTORE_NAMESPACE") // optional; leave blank for default

	ds, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("datastore client: %v", err)
	}
	seed.Seed(ctx, ds, namespace) // idempotent
	_ = ds.Close()

	/* ---------------------------------------------------------
	   2. Normal HTTP server using your existing router
	--------------------------------------------------------- */
	router := setupRouter() // all routes live in router.go

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
