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

type CommentHandler struct {
	service commentService
}

type commentService interface {
	Create(req dto.CreateCommentRequest, offerID string, userID uuid.UUID) error
	FindByOfferID(offerID string) ([]dto.CommentResponse, error)
}

func NewCommentHandler(service commentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}

	offerID := chi.URLParam(r, "id")

	var req dto.CreateCommentRequest

	err := httpx.DecodeJSON(w, r, &req)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json", "The request body must be valid JSON with only supported fields.")
		return
	}

	err = h.service.Create(req, offerID, userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			httpx.WriteError(w, http.StatusBadRequest, "invalid_comment", "The comment must contain between 1 and 1000 characters and the offer id must be valid.")
		case errors.Is(err, services.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, "offer_not_found", "The offer was not found.")
		default:
			log.Printf("create comment: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	offerID := chi.URLParam(r, "id")

	comments, err := h.service.FindByOfferID(offerID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			httpx.WriteError(w, http.StatusBadRequest, "invalid_offer_id", "The offer id is invalid.")
		case errors.Is(err, services.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, "offer_not_found", "The offer was not found.")
		default:
			log.Printf("list comments: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, comments)
}
