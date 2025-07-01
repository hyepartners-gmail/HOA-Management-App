// models/user.go
package models

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/datastore"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
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
	PasswordResetToken     string     `datastore:"password_reset_token,omitempty" json:"-"`
}

func UpdateUser(user *User) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("User", user.ID, nil)
	_, err := client.Put(ctx, key, user)
	return err
}

func SaveUser(user User) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("User", user.ID, nil)
	_, err := client.Put(ctx, key, &user)
	return err
}

func UpdateUserPassword(userID string, hashed string) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("User", userID, nil)
	_, err := client.Mutate(ctx, datastore.NewUpdate(key, map[string]interface{}{
		"hashed_password": hashed,
	}))
	return err
}

func GetUserByEmail(email string) (*User, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	q := datastore.NewQuery("User").Filter("email =", email).Limit(1)
	var users []*User
	_, err := client.GetAll(ctx, q, &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}
