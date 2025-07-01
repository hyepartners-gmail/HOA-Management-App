package models

import (
	"context"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
)

type SpecialAssessment struct {
	ID        uuid.UUID  `datastore:"id"`
	CabinID   uuid.UUID  `datastore:"cabin_id"`
	OwnerID   uuid.UUID  `datastore:"owner_id"`
	Reason    string     `datastore:"reason"`
	Date      time.Time  `datastore:"date"`
	Share     float64    `datastore:"share"`
	AmountDue float64    `datastore:"amount_due"`
	Paid      bool       `datastore:"paid"`
	PaidAt    *time.Time `datastore:"paid_at,omitempty"`
	CreatedAt time.Time  `datastore:"created_at"`
}

func GenerateSpecialAssessments(reason string, date time.Time, total float64) error {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)

	cabins, err := GetAllCabins()
	if err != nil {
		return err
	}

	var assessments []*SpecialAssessment
	for _, cabin := range cabins {
		owners, _ := GetOwnersByCabinID(cabin.ID)
		if len(owners) == 0 {
			continue
		}
		primary := owners[0]
		share := float64(cabin.ShareCount) / 100.0
		amount := total * share

		assessments = append(assessments, &SpecialAssessment{
			ID:        uuid.New(),
			CabinID:   cabin.ID,
			OwnerID:   primary.ID,
			Reason:    reason,
			Date:      date,
			Share:     share,
			AmountDue: amount,
			CreatedAt: time.Now(),
		})
	}

	for _, a := range assessments {
		key := datastore.NameKey("SpecialAssessment", a.ID.String(), nil)
		if _, err := client.Put(ctx, key, a); err != nil {
			return err
		}
	}

	return nil
}

func GetAssessmentsByOwnerID(ownerID uuid.UUID) ([]*SpecialAssessment, error) {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)

	var results []*SpecialAssessment
	query := datastore.NewQuery("SpecialAssessment").Filter("owner_id =", ownerID)
	_, err := client.GetAll(ctx, query, &results)
	return results, err
}
