// models/owner.go
package models

import (
	"context"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
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
	client := utils.GetDatastoreClient(ctx)

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
