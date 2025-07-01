package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/authutils"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"
)

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	owner, err := models.GetOwnerByID(user.AssociatedOwnerID)
	if err != nil {
		utils.JSONError(w, "Owner not found", http.StatusNotFound)
		return
	}
	response := map[string]interface{}{
		"user":  user,
		"owner": owner,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var input struct {
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		MailingAddress string `json:"mailing_address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ownerID, err := uuid.Parse(user.AssociatedOwnerID)
	if err != nil {
		utils.JSONError(w, "Invalid owner ID", http.StatusBadRequest)
		return
	}
	err = models.UpdateOwnerContactInfo(ownerID, input.Email, input.Phone, input.MailingAddress)

	if err != nil {
		utils.JSONError(w, "Failed to update", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var input struct {
		Current string `json:"current"`
		New     string `json:"new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if !authutils.CheckPasswordHash(input.Current, user.HashedPassword) {
		utils.JSONError(w, "Incorrect current password", http.StatusUnauthorized)
		return
	}

	newHash, err := authutils.HashPassword(input.New)
	if err != nil {
		utils.JSONError(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// userID, err := uuid.Parse(user.ID)
	// if err != nil {
	// 	utils.JSONError(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }
	// err = models.UpdateUserPassword(userID, newHash)
	err = models.UpdateUserPassword(user.ID, newHash)

	if err != nil {
		utils.JSONError(w, "Failed to save password", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
