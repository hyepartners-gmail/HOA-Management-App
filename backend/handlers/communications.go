package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/go-chi/chi/v5"
)

func ListCommunicationsHandler(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	limit := 100
	communications, err := models.ListCommunications(limit, t)
	if err != nil {
		utils.JSONError(w, "Failed to load", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(communications)
}

func GetCommunicationHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	comm, err := models.GetCommunicationByID(id)
	if err != nil {
		utils.JSONError(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(comm)
}
