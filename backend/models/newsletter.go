package models

import (
	"time"

	"github.com/google/uuid"
)

type Newsletter struct {
	ID              uuid.UUID  `datastore:"id" json:"id"`
	Title           string     `datastore:"title" json:"title"`
	Body            string     `datastore:"body,noindex" json:"body"` // allow long form
	PublishedAt     *time.Time `datastore:"published_at" json:"published_at,omitempty"`
	CreatedByUserID uuid.UUID  `datastore:"created_by_user_id" json:"created_by_user_id"`
}

func SaveNewsletter(n *Newsletter) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return saveToDatastore("newsletters", n.ID.String(), n)
}

func GetAllNewsletters() ([]Newsletter, error) {
	var newsletters []Newsletter
	err := queryAllFromDatastore("newsletters", &newsletters)
	return newsletters, err
}

func GetNewsletterByID(id string) (*Newsletter, error) {
	var n Newsletter
	err := loadFromDatastore(&n, "newsletters", id)
	return &n, err
}
