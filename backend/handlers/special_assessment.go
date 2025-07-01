package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

type SpecialAssessmentInput struct {
	Reason      string    `json:"reason"`
	Date        time.Time `json:"date"`
	TotalAmount float64   `json:"total_amount"`
}

func TriggerAssessmentHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "admin" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input SpecialAssessmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := models.GenerateSpecialAssessments(input.Reason, input.Date, input.TotalAmount); err != nil {
		utils.JSONError(w, "Failed to create assessments", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetMyAssessmentsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.AssociatedOwnerID == nil {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	assessments, err := models.GetAssessmentsByOwnerID(*user.AssociatedOwnerID)
	if err != nil {
		utils.JSONError(w, "Failed to load assessments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(assessments)
}
