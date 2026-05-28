package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/services"
	"github.com/go-chi/chi/v5"
)

type OfferHandler struct {
	service *services.OfferService
}

func NewOfferHandler(service *services.OfferService) *OfferHandler {
	return &OfferHandler{service: service}
}

func (h *OfferHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOfferRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.service.Create(req)
	if err != nil {
		http.Error(w, "failed to create offer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *OfferHandler) FindNearby(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}

	offers, err := h.service.FindNearby(latitude, longitude)
	if err != nil {
		http.Error(w, "failed to fetch offers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(offers)
}

func (h *OfferHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	offer, err := h.service.FindByID(id)
	if err != nil {
		http.Error(w, "offer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(offer)
}
