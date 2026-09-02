package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			http.Error(w, "name, email and a password of at least 8 characters are required", http.StatusBadRequest)
		case errors.Is(err, services.ErrEmailTaken):
			http.Error(w, "email already registered", http.StatusConflict)
		default:
			http.Error(w, "failed to register", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
