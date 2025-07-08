package seed

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
)

//go:embed seed.json
var seedFile embed.FS

// Seed checks every Kind listed in the JSON; if a kind is empty it loads
// the records from mockdata.json and writes them. 100 % idempotent.
func Seed(ctx context.Context, ds *datastore.Client, namespace string) {
	/* -------- 1. Read the file once -------- */
	raw, err := seedFile.ReadFile("mockdata.json")
	if err != nil {
		log.Printf("[seed] cannot read mockdata.json: %v", err)
		return
	}

	var byKind map[string][]map[string]any
	if err := json.Unmarshal(raw, &byKind); err != nil {
		log.Printf("[seed] bad JSON: %v", err)
		return
	}

	/* -------- 2. Iterate every kind present -------- */
	for kind, records := range byKind {
		if len(records) == 0 {
			log.Printf("[seed] %s has 0 records, skipping", kind)
			continue
		}

		// already populated?
		q := datastore.NewQuery(kind).Namespace(namespace).KeysOnly().Limit(1)
		it := ds.Run(ctx, q)
		_, err := it.Next(nil)
		if err != datastore.Done { // at least one entity exists
			log.Printf("[seed] %s already populated, skipping", kind)
			continue
		}

		keys := make([]*datastore.Key, 0, len(records))
		valid := make([]map[string]any, 0, len(records))

		for i, rec := range records {
			id, ok := rec["id"].(string)
			if !ok || id == "" {
				log.Printf("[seed] %s[%d] missing id, skipping", kind, i)
				continue
			}
			key := datastore.NameKey(kind, id, nil)
			key.Namespace = namespace
			keys = append(keys, key)
			valid = append(valid, rec)
		}

		if len(valid) == 0 {
			log.Printf("[seed] %s: nothing valid to write", kind)
			continue
		}

		if _, err := ds.PutMulti(ctx, keys, valid); err != nil {
			log.Printf("[seed] put %s: %v", kind, err)
			continue
		}
		log.Printf("[seed] %s: wrote %d records", kind, len(valid))
	}

	fmt.Println("[seed] completed")
}
