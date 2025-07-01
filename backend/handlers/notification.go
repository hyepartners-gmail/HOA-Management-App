package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/google/uuid"
)

func CreateNotificationHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var n models.Notification
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		utils.JSONError(w, "invalid input", http.StatusBadRequest)
		return
	}
	n.ID = uuid.New()
	uid, err := uuid.Parse(user.ID)
	if err != nil {
		utils.JSONError(w, "invalid user ID", http.StatusInternalServerError)
		return
	}
	n.CreatedByUserID = uid
	n.CreatedAt = time.Now()

	if err := models.SaveNotification(&n); err != nil {
		utils.JSONError(w, "failed to save", http.StatusInternalServerError)
		return
	}

	if n.Type == models.TypeFlash && (n.DeliveryMethod == models.DeliverySMS || n.DeliveryMethod == models.DeliveryBoth) {
		go utils.SendFlashSMS(n)
	}
	if n.DeliveryMethod == models.DeliveryEmail || n.DeliveryMethod == models.DeliveryBoth {
		go utils.SendNotificationEmail(n)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

func ListNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	all, err := models.GetNotificationsForUser(user)
	if err != nil {
		utils.JSONError(w, "failed to fetch", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(all)
}

func UpdateLastSeenNotificationHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	now := time.Now()
	user.LastSeenNotificationAt = &now
	if err := models.UpdateUser(user); err != nil {
		utils.JSONError(w, "failed to update timestamp", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
