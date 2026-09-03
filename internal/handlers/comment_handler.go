package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/services"
	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
	service *services.CommentService
}

func NewCommentHandler(service *services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	offerID := chi.URLParam(r, "id")

	var req dto.CreateCommentRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.service.Create(req, offerID, userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)
		default:
			http.Error(w, "failed to create comment", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	offerID := chi.URLParam(r, "id")

	comments, err := h.service.FindByOfferID(offerID)
	if err != nil {
		http.Error(w, "failed to fetch comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(comments)
}
