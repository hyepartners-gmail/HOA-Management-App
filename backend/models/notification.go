package models

import (
	"time"

	"github.com/google/uuid"
)

type Audience string
type NotificationType string
type DeliveryMethod string

const (
	AudienceAll        Audience = "all"
	AudienceOwnersOnly Audience = "owners_only"
	AudienceBoardOnly  Audience = "board_only"
	AudienceRole       Audience = "specific_roles"

	TypeNormal NotificationType = "normal"
	TypeFlash  NotificationType = "flash"

	DeliveryEmail DeliveryMethod = "email"
	DeliverySMS   DeliveryMethod = "sms"
	DeliveryBoth  DeliveryMethod = "both"
)

type Notification struct {
	ID              uuid.UUID        `datastore:"id"`
	Title           string           `datastore:"title"`
	Body            string           `datastore:"body,noindex"`
	Audience        Audience         `datastore:"audience"`
	Roles           []string         `datastore:"roles"` // only used if Audience == specific_roles
	Type            NotificationType `datastore:"type"`
	DeliveryMethod  DeliveryMethod   `datastore:"delivery_method"`
	CreatedByUserID uuid.UUID        `datastore:"created_by_user_id"`
	CreatedAt       time.Time        `datastore:"created_at"`
	ExpiresAt       *time.Time       `datastore:"expires_at,omitempty"`
}
