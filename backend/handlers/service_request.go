package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/google/uuid"
)

func SubmitServiceRequestHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var req models.ServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "Invalid body", http.StatusBadRequest)
		return
	}

	req.ID = uuid.New()
	req.SubmittedByUserID = uuid.MustParse(user.ID)
	req.CreatedAt = time.Now()
	req.Status = "open"

	if err := models.SaveServiceRequest(req); err != nil {
		utils.JSONError(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetAllServiceRequestsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "admin" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requests, err := models.GetAllServiceRequests()
	if err != nil {
		utils.JSONError(w, "Fetch error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

func UpdateServiceRequestStatusHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "admin" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.JSONError(w, "Missing ID", http.StatusBadRequest)
		return
	}

	if err := models.UpdateServiceRequestStatus(id, input.Status); err != nil {
		utils.JSONError(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
