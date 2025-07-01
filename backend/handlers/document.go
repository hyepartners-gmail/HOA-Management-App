package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/models"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/utils"

	"github.com/google/uuid"
)

func UploadDocumentHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	title := r.FormValue("title")
	category := r.FormValue("category")
	visibleTo := r.FormValue("visible_to")

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.JSONError(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		utils.JSONError(w, "File read error", http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	docID := uuid.New()
	objectName := "documents/" + docID.String() + ext
	url, err := utils.UploadToGCS(buf, objectName)
	if err != nil {
		utils.JSONError(w, "Upload failed", http.StatusInternalServerError)
		return
	}

	doc := models.Document{
		ID:         docID,
		Title:      title,
		Category:   category,
		URL:        url,
		VisibleTo:  visibleTo,
		UploadedBy: user.ID,
		UploadedAt: time.Now(),
	}
	if err := models.SaveDocument(doc); err != nil {
		utils.JSONError(w, "Save failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

func ListDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	docs, err := models.ListDocuments(string(user.Role))
	if err != nil {
		utils.JSONError(w, "Load failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(docs)
}
