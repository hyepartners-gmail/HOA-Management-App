package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

func GetNotificationsForUser(user *User) ([]Notification, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	now := time.Now()
	query := datastore.NewQuery("Notification").
		Filter("expires_at >", now).
		Order("-created_at")

	var results []Notification
	_, err := client.GetAll(ctx, query, &results)
	if err != nil {
		return nil, err
	}

	filtered := []Notification{}
	for _, n := range results {
		switch n.Audience {
		case AudienceAll:
			filtered = append(filtered, n)
		case AudienceOwnersOnly:
			if user.Role == RoleCabinOwner {
				filtered = append(filtered, n)
			}
		case AudienceBoardOnly:
			if user.Role != RoleCabinOwner {
				filtered = append(filtered, n)
			}
		case AudienceRole:
			for _, role := range n.Roles {
				if string(user.Role) == role {
					filtered = append(filtered, n)
					break
				}
			}
		}
	}

	return filtered, nil
}

func ResolveRecipients(n Notification) struct{ Emails, Phones []string } {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var users []User
	query := datastore.NewQuery("User")

	switch n.Audience {
	case AudienceAll:
		// No filter
	case AudienceOwnersOnly:
		query = query.Filter("role =", RoleCabinOwner)
	case AudienceBoardOnly:
		query = query.Filter("role >", RoleCabinOwner) // assumes enum order
	case AudienceRole:
		// manual filter after fetch
	}

	_, err := client.GetAll(ctx, query, &users)
	if err != nil {
		return struct{ Emails, Phones []string }{}
	}

	emails := []string{}
	phones := []string{}
	for _, u := range users {
		if n.Audience == AudienceRole {
			match := false
			for _, r := range n.Roles {
				if string(u.Role) == r {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		emails = append(emails, u.Email)
		// Phone lookup can be added here
	}

	return struct{ Emails, Phones []string }{
		Emails: emails,
		Phones: phones,
	}
}
