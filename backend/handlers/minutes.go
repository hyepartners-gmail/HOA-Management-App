package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

func UploadMeetingMinutesHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "secretary" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var minutes models.MeetingMinutes
	if err := json.NewDecoder(r.Body).Decode(&minutes); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	minutes.UploadedByUserID = user.ID

	if err := models.SaveMeetingMinutes(minutes); err != nil {
		utils.JSONError(w, "Failed to save meeting minutes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetMeetingMinutesHandler(w http.ResponseWriter, r *http.Request) {
	minutes, err := models.GetAllMeetingMinutes()
	if err != nil {
		utils.JSONError(w, "Failed to fetch meeting minutes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(minutes)
}
