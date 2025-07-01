package models

import (
	"time"

	"github.com/google/uuid"
)

type PostCategory string

const (
	ForSale            PostCategory = "for_sale"
	ForFree            PostCategory = "for_free"
	Wanted             PostCategory = "wanted"
	CommunityStories   PostCategory = "community_stories"
	AnimalSightings    PostCategory = "animal_sightings"
	FlashNotifications PostCategory = "flash_notifications"
	CanYouHelpMe       PostCategory = "can_you_help_me"
	AfterTheStorm      PostCategory = "after_the_storm"
	Rideshare          PostCategory = "rideshare"
)

type Post struct {
	ID            uuid.UUID    `datastore:"id"`
	Category      PostCategory `datastore:"category"`
	Title         string       `datastore:"title"`
	Body          string       `datastore:"body,noindex"`
	CreatedByUser uuid.UUID    `datastore:"created_by_user_id"`
	CreatedAt     time.Time    `datastore:"created_at"`
}
