// models/user.go
package models

import (
	"backend/utils"
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin      Role = "admin"
	RolePresident  Role = "president"
	RoleSecretary  Role = "secretary"
	RoleTreasurer  Role = "treasurer"
	RoleCabinOwner Role = "cabin_owner"
)

type User struct {
	ID                     string     `datastore:"id"`
	Email                  string     `datastore:"email"`
	HashedPassword         string     `datastore:"hashed_password"`
	Role                   Role       `datastore:"role"`
	AssociatedOwnerID      string     `datastore:"associated_owner_id"`
	LastSeenNotificationAt *time.Time `datastore:"last_seen_notification_at"`
}

func UpdateUserPassword(userID uuid.UUID, newHash string) error {
	ctx := context.Background()
	client := utils.GetDatastoreClient(ctx)

	var user User
	key := datastore.NameKey("User", userID.String(), nil)
	if err := client.Get(ctx, key, &user); err != nil {
		return err
	}
	user.HashedPassword = newHash

	_, err := client.Put(ctx, key, &user)
	return err
}
