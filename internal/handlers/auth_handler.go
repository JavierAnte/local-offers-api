package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/httpx"
	"github.com/JavierAnte/local-offers-api/internal/services"
)

type AuthHandler struct {
	service authService
}

type authService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

func NewAuthHandler(service authService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json", "The request body must be valid JSON with only supported fields.")
		return
	}

	resp, err := h.service.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidInput):
			httpx.WriteError(w, http.StatusBadRequest, "invalid_registration", "Name, a valid email, and a password of 8 to 72 characters are required.")
		case errors.Is(err, services.ErrEmailTaken):
			httpx.WriteError(w, http.StatusConflict, "email_taken", "The email is already registered.")
		default:
			log.Printf("register user: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_json", "The request body must be valid JSON with only supported fields.")
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "The email or password is incorrect.")
			return
		}
		log.Printf("login user: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}
