package models

func SavePost(p *Post) error {
	return saveToDatastore(p, "posts", p.ID.String())
}

func GetAllPosts(category string) ([]Post, error) {
	// Optional: filter by category
	// Otherwise sort by CreatedAt descending
}

func SaveComment(c *Comment) error {
	return saveToDatastore(c, "comments", c.ID.String())
}

func GetCommentsForPost(postID string) ([]Comment, error) {
	// Query comments where post_id = postID, sort by CreatedAt asc
}

func GetPostByID(id string) (*Post, error) {
	var post Post
	err := loadFromDatastore(&post, "posts", id)
	return &post, err
}

func DeletePost(id string) error {
	return deleteFromDatastore("posts", id)
}

func DeleteCommentsForPost(postID string) error {
	// Query comments by post_id and delete each
	// This is optional cleanup
	return deleteByField("comments", "post_id", postID)
}
