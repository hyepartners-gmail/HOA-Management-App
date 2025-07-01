package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.JSONError(w, "Invalid request", http.StatusBadRequest)

		return
	}

	user, err := models.AuthenticateUser(creds.Email, creds.Password)
	if err != nil {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {

		utils.JSONError(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func PasswordResetRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for request email logic
	w.WriteHeader(http.StatusOK)
}

func PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for reset logic
	w.WriteHeader(http.StatusOK)
}
