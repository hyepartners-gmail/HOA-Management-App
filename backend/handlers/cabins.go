package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/go-chi/chi/v5"
)

func GetCabins(w http.ResponseWriter, r *http.Request) {
	cabins, err := models.GetAllCabins()
	if err != nil {
		utils.JSONError(w, "Failed to load cabins", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cabins)
}

func CreateCabin(w http.ResponseWriter, r *http.Request) {
	var c models.Cabin
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.JSONError(w, "Bad request", http.StatusBadRequest)
		return
	}
	models.SaveCabin(&c)
	w.WriteHeader(http.StatusCreated)
}

func UpdateCabin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var c models.Cabin
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.JSONError(w, "Bad request", http.StatusBadRequest)
		return
	}
	c.ID = id
	models.UpdateCabin(&c)
	w.WriteHeader(http.StatusOK)
}

func GetCabinByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cabin, err := models.FindCabinByID(id)
	if err != nil {
		utils.JSONError(w, "Cabin not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(cabin)
}
