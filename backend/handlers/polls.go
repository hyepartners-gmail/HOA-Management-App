package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/go-chi/chi/v5"
)

func CreatePollHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "president" && user.Role != "secretary" && user.Role != "treasurer" && user.Role != "admin" {
		utils.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		Question  string    `json:"question"`
		Options   []string  `json:"options"`
		Audience  string    `json:"audience"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.JSONError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(payload.Options) < 2 || len(payload.Options) > 6 {
		utils.JSONError(w, "poll must have 2–6 options", http.StatusBadRequest)
		return
	}

	p := models.Poll{
		Question:  payload.Question,
		Options:   payload.Options,
		Audience:  payload.Audience,
		StartDate: payload.StartDate,
		EndDate:   payload.EndDate,
		CreatedBy: user.ID,
	}

	if err := models.CreatePoll(p); err != nil {
		utils.JSONError(w, "failed to create poll", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ListPollsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	role := "all"
	if user.Role == "cabin_owner" {
		role = "owners"
	} else if user.Role != "admin" {
		role = "board"
	}

	polls, err := models.GetPollsForUser(role)
	if err != nil {
		utils.JSONError(w, "load failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(polls)
}

func SubmitVoteHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	pollID := chi.URLParam(r, "id")

	var body struct {
		Choice int `json:"choice"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.JSONError(w, "invalid vote", http.StatusBadRequest)
		return
	}

	voted, err := models.HasVoted(pollID, user.ID)
	if err != nil {
		utils.JSONError(w, "error checking vote", http.StatusInternalServerError)
		return
	}
	if voted {
		utils.JSONError(w, "already voted", http.StatusForbidden)
		return
	}

	vote := models.Vote{
		PollID: pollID,
		UserID: user.ID,
		Choice: body.Choice,
	}

	if err := models.SubmitVote(vote); err != nil {
		utils.JSONError(w, "vote failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
