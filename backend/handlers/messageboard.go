package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category") // optional
	posts, err := models.GetAllPosts(category)
	if err != nil {
		utils.JSONError(w, "failed to load posts", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.JSONError(w, "invalid payload", http.StatusBadRequest)
		return
	}
	p.ID = uuid.New()
	p.CreatedAt = time.Now()
	p.CreatedByUser = user.ID

	if err := models.SavePost(&p); err != nil {
		utils.JSONError(w, "failed to save post", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func ListCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	comments, err := models.GetCommentsForPost(postID)
	if err != nil {
		utils.JSONError(w, "failed to load comments", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(comments)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	postID := chi.URLParam(r, "postID")

	var c models.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.JSONError(w, "invalid comment", http.StatusBadRequest)
		return
	}
	c.ID = uuid.New()
	c.PostID = uuid.MustParse(postID)
	c.UserID = user.ID
	c.CreatedAt = time.Now()

	if err := models.SaveComment(&c); err != nil {
		utils.JSONError(w, "failed to save comment", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(c)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	postID := chi.URLParam(r, "postID")

	post, err := models.GetPostByID(postID)
	if err != nil {
		utils.JSONError(w, "post not found", http.StatusNotFound)
		return
	}

	if post.CreatedByUser != user.ID {
		utils.JSONError(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := models.DeletePost(postID); err != nil {
		utils.JSONError(w, "failed to delete post", http.StatusInternalServerError)
		return
	}

	// Optionally: delete comments on this post as well
	_ = models.DeleteCommentsForPost(postID)

	w.WriteHeader(http.StatusNoContent)
}
