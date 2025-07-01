// models/owner.go
package models

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type Owner struct {
	ID             string   `datastore:"id"`
	FullName       string   `datastore:"full_name"`
	Email          string   `datastore:"email"`
	Phone          string   `datastore:"phone"`
	MailingAddress string   `datastore:"mailing_address"`
	CabinIDs       []string `datastore:"cabin_ids"`
	IsPrimary      bool     `datastore:"is_primary"`
	LoginEnabled   bool     `datastore:"login_enabled"`
}

func UpdateOwnerContactInfo(ownerID uuid.UUID, email, phone, mailingAddr string) error {
	ctx := context.Background()
	// client := utils.GetDatastoreClient(ctx)
	client := ds.GetClient(ctx)

	var owner Owner
	key := datastore.NameKey("Owner", ownerID.String(), nil)
	if err := client.Get(ctx, key, &owner); err != nil {
		return err
	}

	owner.Email = email
	owner.Phone = phone
	owner.MailingAddress = mailingAddr

	_, err := client.Put(ctx, key, &owner)
	return err
}

func GetOwnerByID(id string) (*Owner, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)
	key := datastore.NameKey("Owner", id, nil)
	var o Owner
	if err := client.Get(ctx, key, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func GetOwnersByCabinID(cabinID string) ([]*Owner, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var owners []*Owner
	query := datastore.NewQuery("Owner").Filter("cabin_id =", cabinID)
	_, err := client.GetAll(ctx, query, &owners)
	return owners, err
}

func GetAllOwners() ([]*Owner, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var owners []*Owner
	query := datastore.NewQuery("Owner")
	_, err := client.GetAll(ctx, query, &owners)
	return owners, err
}

func CreateOwner(
	id uuid.UUID,
	fullName, email, phone, mailingAddress string,
	cabinIDs []string,
	isPrimary, loginEnabled bool,
) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	owner := Owner{
		ID:             id.String(),
		FullName:       fullName,
		Email:          email,
		Phone:          phone,
		MailingAddress: mailingAddress,
		CabinIDs:       cabinIDs,
		IsPrimary:      isPrimary,
		LoginEnabled:   loginEnabled,
	}

	key := datastore.NameKey("Owner", id.String(), nil)
	_, err := client.Put(ctx, key, &owner)
	return err
}

func UpdateOwner(owner *Owner) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("Owner", owner.ID, nil)
	_, err := client.Put(ctx, key, owner)
	return err
}
