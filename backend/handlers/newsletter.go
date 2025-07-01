package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

func ListNewslettersHandler(w http.ResponseWriter, r *http.Request) {
	newsletters, err := models.GetAllNewsletters()
	if err != nil {
		utils.JSONError(w, "could not load newsletters", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newsletters)
}

func CreateNewsletterHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var payload struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.JSONError(w, "invalid input", http.StatusBadRequest)
		return
	}

	newsletter := &models.Newsletter{
		ID:              uuid.New(),
		Title:           payload.Title,
		Body:            payload.Body,
		CreatedByUserID: user.ID,
	}

	if err := models.SaveNewsletter(newsletter); err != nil {
		utils.JSONError(w, "failed to save newsletter", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newsletter)
}

func PublishNewsletterHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(*models.User)

	n, err := models.GetNewsletterByID(id)
	if err != nil {
		utils.JSONError(w, "not found", http.StatusNotFound)
		return
	}

	now := time.Now()
	n.PublishedAt = &now

	if err := models.SaveNewsletter(n); err != nil {
		utils.JSONError(w, "failed to publish", http.StatusInternalServerError)
		return
	}

	// Send email to all cabin owners
	go utils.SendNewsletterToAllOwners(n.Title, n.Body)

	w.WriteHeader(http.StatusOK)
}
