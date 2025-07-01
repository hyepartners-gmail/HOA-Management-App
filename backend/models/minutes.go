package models

import "time"

type MeetingMinutes struct {
	ID          string    `datastore:"id"`
	MeetingDate time.Time `datastore:"meeting_date"`
	Title       string    `datastore:"title"`
	ContentHTML string    `datastore:"content_html"` // or URL to uploaded PDF
	Published   bool      `datastore:"published"`
	CreatedBy   string    `datastore:"created_by_user_id"`
	CreatedAt   time.Time `datastore:"created_at"`
}
