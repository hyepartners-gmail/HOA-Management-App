package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
)

type Poll struct {
	ID        uuid.UUID `datastore:"id" json:"id"`
	Question  string    `datastore:"question" json:"question"`
	Options   []string  `datastore:"options" json:"options"`
	Audience  string    `datastore:"audience" json:"audience"` // e.g., "all", "owners", "board"
	StartDate time.Time `datastore:"start_date" json:"start_date"`
	EndDate   time.Time `datastore:"end_date" json:"end_date"`
	CreatedBy uuid.UUID `datastore:"created_by" json:"created_by"`
}

type Vote struct {
	PollID      string    `datastore:"poll_id" json:"poll_id"`
	UserID      uuid.UUID `datastore:"user_id" json:"user_id"`
	Choice      int       `datastore:"choice" json:"choice"`
	SubmittedAt time.Time `datastore:"submitted_at" json:"submitted_at"`
}

func CreatePoll(p Poll) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	p.ID = uuid.New()
	key := datastore.NameKey("Poll", p.ID.String(), nil)
	_, err := client.Put(ctx, key, &p)
	return err
}

func SubmitVote(v Vote) error {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	key := datastore.NameKey("Vote", v.PollID+"_"+v.UserID.String(), nil)
	v.SubmittedAt = time.Now()
	_, err := client.Put(ctx, key, &v)
	return err
}

func HasVoted(pollID string, userID uuid.UUID) (bool, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	var v Vote
	key := datastore.NameKey("Vote", pollID+"_"+userID.String(), nil)
	err := client.Get(ctx, key, &v)
	if err == datastore.ErrNoSuchEntity {
		return false, nil
	}
	return err == nil, err
}

func GetPollsForUser(audience string) ([]Poll, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	now := time.Now()
	q := datastore.NewQuery("Poll").
		Filter("Audience =", audience).
		Filter("StartDate <=", now).
		Filter("EndDate >=", now).
		Order("-StartDate")

	var results []Poll
	_, err := client.GetAll(ctx, q, &results)
	return results, err
}

func GetVotesByPollID(pollID string) ([]Vote, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	q := datastore.NewQuery("Vote").Filter("PollID =", pollID)
	var results []Vote
	_, err := client.GetAll(ctx, q, &results)
	return results, err
}
