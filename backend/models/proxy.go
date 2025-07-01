package models

import "time"

type ProxyAssignment struct {
	ID          string     `datastore:"id"`
	FromUserID  string     `datastore:"from_user_id"`
	ToUserID    *string    `datastore:"to_user_id,omitempty"` // optional if to_office used
	ToOffice    *string    `datastore:"to_office,omitempty"`  // enum: president, secretary, etc
	MeetingDate *time.Time `datastore:"meeting_date,omitempty"`
	IsOneTime   bool       `datastore:"is_one_time"`
	CreatedAt   time.Time  `datastore:"created_at"`
}
