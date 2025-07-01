package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `datastore:"id"`
	PostID    uuid.UUID `datastore:"post_id"`
	UserID    uuid.UUID `datastore:"user_id"`
	Body      string    `datastore:"body,noindex"`
	CreatedAt time.Time `datastore:"created_at"`
}
