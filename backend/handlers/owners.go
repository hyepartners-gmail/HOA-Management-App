package handlers

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetOwners(w http.ResponseWriter, r *http.Request) {
	owners := models.GetAllOwners()
	json.NewEncoder(w).Encode(owners)
}

func CreateOwner(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName       string   `json:"full_name"`
		Email          string   `json:"email"`
		Phone          string   `json:"phone"`
		MailingAddress string   `json:"mailing_address"`
		CabinIDs       []string `json:"cabin_ids"`
		IsPrimary      bool     `json:"is_primary"`
		LoginEnabled   bool     `json:"login_enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ownerID := uuid.New()
	err := models.CreateOwner(ownerID, input.FullName, input.Email, input.Phone, input.MailingAddress, input.CabinIDs, input.IsPrimary, input.LoginEnabled)
	if err != nil {
		utils.JSONError(w, "Failed to create owner", http.StatusInternalServerError)
		return
	}

	// Create associated User if login enabled
	if input.LoginEnabled {
		userID := uuid.New()
		passwordResetToken := utils.GenerateResetToken() // assume secure token gen
		user := models.User{
			ID:                 userID,
			Email:              input.Email,
			HashedPassword:     "", // blank for now
			Role:               "cabin_owner",
			AssociatedOwnerID:  ownerID,
			PasswordResetToken: passwordResetToken,
		}
		if err := models.SaveUser(user); err != nil {
			utils.JSONError(w, "User creation failed", http.StatusInternalServerError)
			return
		}

		// Send email with reset token
		err = utils.SendWelcomeEmail(input.Email, passwordResetToken)
		if err != nil {
			utils.JSONError(w, "Email failed", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateOwner(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var o models.Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		utils.JSONError(w, "Bad request", http.StatusBadRequest)
		return
	}
	o.ID = id
	models.UpdateOwner(&o)
	w.WriteHeader(http.StatusOK)
}
