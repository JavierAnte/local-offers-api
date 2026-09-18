package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/httpx"
	"github.com/JavierAnte/local-offers-api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OfferVoteHandler struct {
	service offerVoteService
}

type offerVoteService interface {
	Vote(offerID string, userID uuid.UUID, voteType string) (dto.VoteResponse, error)
}

func NewOfferVoteHandler(service offerVoteService) *OfferVoteHandler {
	return &OfferVoteHandler{service: service}
}

func (h *OfferVoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}

	offerID := chi.URLParam(r, "id")

	var req dto.CreateVoteRequest

	err := httpx.DecodeJSON(w, r, &req)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json", "The request body must be valid JSON with only supported fields.")
		return
	}

	result, err := h.service.Vote(offerID, userID, req.Type)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			httpx.WriteError(w, http.StatusBadRequest, "invalid_vote", "The offer id and vote type must be valid.")
		case errors.Is(err, services.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, "offer_not_found", "The offer was not found.")
		default:
			log.Printf("register vote: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, result)
}
