package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/google/uuid"
)

func SubmitTalentHandler(w http.ResponseWriter, r *http.Request) {
	var input models.TalentListing
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	input.ID = uuid.New()
	input.CreatedAt = time.Now()
	input.IsApproved = true // auto-list unless flagged
	if err := models.SaveTalentListing(input); err != nil {
		utils.JSONError(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetPublicTalentHandler(w http.ResponseWriter, r *http.Request) {
	listings, err := models.GetApprovedTalentListings()
	if err != nil {
		utils.JSONError(w, "Failed to load", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(listings)
}

func GetAllTalentHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "admin" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	listings, err := models.GetAllTalentListings()
	if err != nil {
		utils.JSONError(w, "Failed to load", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(listings)
}

func ToggleTalentApprovalHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "admin" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		ID       string `json:"id"`
		Approved bool   `json:"approved"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := models.ApproveTalentListing(input.ID, input.Approved); err != nil {
		utils.JSONError(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
