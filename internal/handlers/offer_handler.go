package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/httpx"
	"github.com/JavierAnte/local-offers-api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OfferHandler struct {
	service offerService
}

type offerService interface {
	Create(req dto.CreateOfferRequest, userID uuid.UUID) error
	FindNearby(latitude float64, longitude float64, radiusMeters int) ([]dto.OfferResponse, error)
	FindByID(id string) (*dto.OfferResponse, error)
}

func NewOfferHandler(service offerService) *OfferHandler {
	return &OfferHandler{service: service}
}

func (h *OfferHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}

	var req dto.CreateOfferRequest

	err := httpx.DecodeJSON(w, r, &req)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json", "The request body must be valid JSON with only supported fields.")
		return
	}

	err = h.service.Create(req, userID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			httpx.WriteError(w, http.StatusBadRequest, "invalid_offer", "One or more offer fields are invalid.")
			return
		}
		log.Printf("create offer: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

const (
	defaultRadiusMeters = 5000
	maxRadiusMeters     = 20000
)

func (h *OfferHandler) FindNearby(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_coordinates", "Latitude and longitude must be valid coordinates.")
		return
	}

	longitude, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_coordinates", "Latitude and longitude must be valid coordinates.")
		return
	}

	radius := defaultRadiusMeters
	if radiusStr := r.URL.Query().Get("radius"); radiusStr != "" {
		parsed, err := strconv.Atoi(radiusStr)
		if err != nil || parsed <= 0 {
			httpx.WriteError(w, http.StatusBadRequest, "invalid_radius", "Radius must be a positive integer.")
			return
		}
		radius = parsed
		if radius > maxRadiusMeters {
			radius = maxRadiusMeters
		}
	}

	offers, err := h.service.FindNearby(latitude, longitude, radius)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			httpx.WriteError(w, http.StatusBadRequest, "invalid_coordinates", "Latitude and longitude must be valid coordinates.")
			return
		}
		log.Printf("find nearby offers: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, offers)
}

func (h *OfferHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	offer, err := h.service.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			httpx.WriteError(w, http.StatusBadRequest, "invalid_offer_id", "The offer id is invalid.")
		case errors.Is(err, services.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, "offer_not_found", "The offer was not found.")
		default:
			log.Printf("find offer: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, offer)
}
