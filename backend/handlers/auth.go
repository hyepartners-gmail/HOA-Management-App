package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/authutils"
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
	token, err := utils.GenerateJWT(user.ID, string(user.Role))

	if err != nil {

		utils.JSONError(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func PasswordResetRequestHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByEmail(input.Email)
	if err != nil {
		utils.JSONError(w, "Email not found", http.StatusNotFound)
		return
	}

	token := uuid.New().String()
	if err := models.StoreResetToken(user.ID, token, time.Now().Add(1*time.Hour)); err != nil {
		utils.JSONError(w, "Could not generate reset token", http.StatusInternalServerError)
		return
	}

	// TODO: Actually send the email
	// sendResetEmail(user.Email, token)

	w.WriteHeader(http.StatusOK)
}

func PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID, err := models.ValidateResetToken(input.Token)
	if err != nil {
		utils.JSONError(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	hashed, err := authutils.HashPassword(input.NewPassword)
	if err != nil {
		utils.JSONError(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	if err := models.UpdateUserPassword(userID, hashed); err != nil {
		utils.JSONError(w, "Could not update password", http.StatusInternalServerError)
		return
	}

	models.DeleteResetToken(input.Token)

	w.WriteHeader(http.StatusOK)
}
