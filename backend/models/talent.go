package models

import (
	"context"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
)

type TalentCategory string

const (
	Realtor     TalentCategory = "Realtor"
	Hauling     TalentCategory = "Hauling"
	AirbnbMgmt  TalentCategory = "AirBNB Management"
	Cleaning    TalentCategory = "Cleaning"
	Plumber     TalentCategory = "Plumber"
	Electrician TalentCategory = "Electrician"
	Contractor  TalentCategory = "Contractor"
	Other       TalentCategory = "Other"
)

type TalentListing struct {
	ID          uuid.UUID      `datastore:"id" json:"id"`
	Name        string         `datastore:"name" json:"name"`
	Category    TalentCategory `datastore:"category" json:"category"`
	SubmittedBy string         `datastore:"submitted_by" json:"submitted_by"`
	Phone       string         `datastore:"phone" json:"phone"`
	Email       string         `datastore:"email" json:"email"`
	Website     string         `datastore:"website" json:"website,omitempty"`
	Description string         `datastore:"description" json:"description,omitempty"`
	IsApproved  bool           `datastore:"is_approved" json:"is_approved"`
	CreatedAt   time.Time      `datastore:"created_at" json:"created_at"`
}

func SaveTalentListing(listing TalentListing) error {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)
	key := datastore.NameKey("TalentListing", listing.ID.String(), nil)
	_, err := client.Put(ctx, key, &listing)
	return err
}

func GetApprovedTalentListings() ([]*TalentListing, error) {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)
	var listings []*TalentListing
	query := datastore.NewQuery("TalentListing").Filter("is_approved =", true).Order("category").Order("name")
	_, err := client.GetAll(ctx, query, &listings)
	return listings, err
}

func GetAllTalentListings() ([]*TalentListing, error) {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)
	var listings []*TalentListing
	query := datastore.NewQuery("TalentListing").Order("-created_at")
	_, err := client.GetAll(ctx, query, &listings)
	return listings, err
}

func ApproveTalentListing(id string, approved bool) error {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)
	key := datastore.NameKey("TalentListing", id, nil)

	var listing TalentListing
	if err := client.Get(ctx, key, &listing); err != nil {
		return err
	}
	listing.IsApproved = approved
	_, err := client.Put(ctx, key, &listing)
	return err
}
