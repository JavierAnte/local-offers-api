package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/services"
)

// maxUploadBytes bounds the request body before it's even parsed. The app
// compresses images to a few hundred KB before upload, so this is a
// generous ceiling meant to guard against abuse, not a target size.
const maxUploadBytes = 10 << 20 // 10MB

type UploadHandler struct {
	service *services.UploadService
}

func NewUploadHandler(service *services.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

func (h *UploadHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, "invalid or too large upload (max 10MB)", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := h.service.Save(file)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			http.Error(w, "invalid image", http.StatusBadRequest)
		default:
			http.Error(w, "failed to save image", http.StatusInternalServerError)
		}
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	result := dto.UploadResponse{
		URL: fmt.Sprintf("%s://%s/uploads/%s", scheme, r.Host, filename),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(result)
}
