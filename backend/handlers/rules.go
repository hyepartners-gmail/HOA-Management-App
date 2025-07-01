package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

func GetHOARulesHandler(w http.ResponseWriter, r *http.Request) {
	rules, err := models.GetHOARules()
	if err != nil {
		utils.JSONError(w, "Failed to fetch HOA rules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

func UpdateHOARulesHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "admin" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updated models.HOARules
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated.LastUpdatedBy = user.ID.String()

	if err := models.SaveHOARules(updated); err != nil {
		utils.JSONError(w, "Failed to save HOA rules", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
