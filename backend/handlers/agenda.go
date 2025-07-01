package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

func SubmitAgendaRequestHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var req models.AgendaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.SubmittedByUserID = user.ID

	if err := models.SaveAgendaRequest(req); err != nil {
		utils.JSONError(w, "Failed to save agenda request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetAgendaRequestsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "secretary" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requests, err := models.GetAllAgendaRequests()
	if err != nil {
		utils.JSONError(w, "Failed to retrieve requests", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}
