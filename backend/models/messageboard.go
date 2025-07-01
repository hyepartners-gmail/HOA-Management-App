package models

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	ds "github.com/hyepartners-gmail/HOA-Management-App/backend/datastore"
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
	ID            uuid.UUID    `datastore:"id" json:"id"`
	Category      PostCategory `datastore:"category" json:"category"`
	Title         string       `datastore:"title" json:"title"`
	Body          string       `datastore:"body,noindex" json:"body"`
	CreatedByUser uuid.UUID    `datastore:"created_by_user_id" json:"created_by_user"`
	CreatedAt     time.Time    `datastore:"created_at" json:"created_at"`
}

type Comment struct {
	ID        uuid.UUID `datastore:"id" json:"id"`
	PostID    uuid.UUID `datastore:"post_id" json:"post_id"`
	UserID    uuid.UUID `datastore:"user_id" json:"user_id"`
	Content   string    `datastore:"content" json:"content"`
	CreatedAt time.Time `datastore:"created_at" json:"created_at"`
}

func GetAllPosts(category string) ([]*Post, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	query := datastore.NewQuery("Post").Order("-created_at")
	if category != "" {
		query = query.Filter("category =", category)
	}

	var posts []*Post
	_, err := client.GetAll(ctx, query, &posts)
	return posts, err
}

func SavePost(p *Post) error {
	return saveToDatastore(p, "Post", p.ID.String())
}

func GetCommentsForPost(postID string) ([]*Comment, error) {
	ctx := context.Background()
	client := ds.GetClient(ctx)

	query := datastore.NewQuery("Comment").Filter("post_id =", postID).Order("created_at")
	var comments []*Comment
	_, err := client.GetAll(ctx, query, &comments)
	return comments, err
}

func SaveComment(c *Comment) error {
	return saveToDatastore(c, "Comment", c.ID.String())
}

func GetPostByID(id string) (*Post, error) {
	var post Post
	err := loadFromDatastore(&post, "Post", id)
	return &post, err
}

func DeletePost(id string) error {
	return deleteFromDatastore("Post", id)
}

func DeleteCommentsForPost(postID string) error {
	return deleteByField("Comment", "post_id", postID)
}
