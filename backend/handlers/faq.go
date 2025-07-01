package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/google/uuid"
)

func GetFAQsHandler(w http.ResponseWriter, r *http.Request) {
	faqs, err := models.GetAllFAQs()
	if err != nil {
		utils.JSONError(w, "Failed to load FAQs", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(faqs)
}

func SaveFAQHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "admin" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.FAQ
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.ID == uuid.Nil {
		input.ID = uuid.New()
		input.CreatedAt = time.Now()
	}
	input.UpdatedAt = time.Now()

	if err := models.SaveFAQ(input); err != nil {
		utils.JSONError(w, "Failed to save", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
