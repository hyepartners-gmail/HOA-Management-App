package models

import "time"

type AgendaRequest struct {
	ID                string    `datastore:"id"`
	SubmittedByUserID string    `datastore:"submitted_by"`
	Subject           string    `datastore:"subject"`
	Description       string    `datastore:"description"`
	RequestedDate     time.Time `datastore:"requested_meeting_date"`
	CreatedAt         time.Time `datastore:"created_at"`
}
