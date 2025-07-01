package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

func ListAuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	if user.Role != "admin" {
		utils.JSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	logs, err := models.ListAuditLogs(200)
	if err != nil {
		utils.JSONError(w, "failed to load audit logs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(logs)
}
