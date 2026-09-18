package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/httpx"
	"github.com/JavierAnte/local-offers-api/internal/services"
)

// maxUploadBytes bounds the request body before it's even parsed. The app
// compresses images to a few hundred KB before upload, so this is a
// generous ceiling meant to guard against abuse, not a target size.
const maxUploadBytes = 10 << 20 // 10MB

type UploadHandler struct {
	service uploadService
}

type uploadService interface {
	Save(file io.Reader) (string, error)
}

func NewUploadHandler(service uploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

func (h *UploadHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_upload", "The upload must be a multipart image no larger than 10 MB.")
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "missing_image", "The multipart field 'image' is required.")
		return
	}
	defer file.Close()

	filename, err := h.service.Save(file)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			httpx.WriteError(w, http.StatusBadRequest, "invalid_image", "The uploaded file is not a supported image.")
		default:
			log.Printf("save upload: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
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

	httpx.WriteJSON(w, http.StatusCreated, result)
}
