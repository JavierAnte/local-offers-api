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

type OfferVoteHandler struct {
	service *services.OfferVoteService
}

func NewOfferVoteHandler(service *services.OfferVoteService) *OfferVoteHandler {
	return &OfferVoteHandler{service: service}
}

func (h *OfferVoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	offerID := chi.URLParam(r, "id")

	var req dto.CreateVoteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Vote(offerID, userID, req.Type)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)
		default:
			http.Error(w, "failed to register vote", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}
