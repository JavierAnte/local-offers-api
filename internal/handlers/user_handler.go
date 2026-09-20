package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/httpx"
	"github.com/JavierAnte/local-offers-api/internal/services"
	"github.com/google/uuid"
)

type UserHandler struct {
	service userService
}

type userService interface {
	Me(userID uuid.UUID) (*dto.UserResponse, error)
	MyOffers(userID uuid.UUID) ([]dto.OfferResponse, error)
}

func (h *UserHandler) MyOffers(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}

	offers, err := h.service.MyOffers(userID)
	if err != nil {
		log.Printf("get current user's offers: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, offers)
}

func NewUserHandler(service userService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}

	user, err := h.service.Me(userID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "user_not_found", "The user was not found.")
			return
		}
		log.Printf("get current user: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, user)
}
