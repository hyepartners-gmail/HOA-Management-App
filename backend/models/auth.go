package models

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/authutils"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthenticateUser(email, password string) (*User, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var users []*User
	query := datastore.NewQuery("User").Filter("email =", email).Limit(1)
	_, err := client.GetAll(ctx, query, &users)
	if err != nil || len(users) == 0 {
		return nil, errors.New("user not found")
	}
	user := users[0]

	if !authutils.CheckPasswordHash(password, user.HashedPassword) {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

var resetTokens = make(map[string]struct {
	UserID    string
	ExpiresAt time.Time
})

func StoreResetToken(userID, token string, expiresAt time.Time) error {
	resetTokens[token] = struct {
		UserID    string
		ExpiresAt time.Time
	}{userID, expiresAt}
	return nil
}

func ValidateResetToken(token string) (string, error) {
	entry, ok := resetTokens[token]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return "", errors.New("invalid or expired token")
	}
	return entry.UserID, nil
}

func DeleteResetToken(token string) {
	delete(resetTokens, token)
}
