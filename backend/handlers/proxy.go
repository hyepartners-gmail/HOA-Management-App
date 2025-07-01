package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

func SubmitProxyHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var proxy models.ProxyAssignment
	if err := json.NewDecoder(r.Body).Decode(&proxy); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	proxy.FromUserID = user.ID

	if err := models.SaveProxy(proxy); err != nil {
		utils.JSONError(w, "Failed to save proxy assignment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetProxiesHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	// Only secretary and admin can view all
	if user.Role != "secretary" && user.Role != "admin" {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	proxies, err := models.GetAllProxies()
	if err != nil {
		utils.JSONError(w, "Failed to retrieve proxies", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(proxies)
}
