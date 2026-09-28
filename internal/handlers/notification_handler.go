package handlers

import (
	"log"
	"net/http"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/httpx"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service notificationService
}

type notificationService interface {
	FindByRecipientID(recipientID uuid.UUID) ([]dto.NotificationResponse, error)
}

func NewNotificationHandler(service notificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication_required", "Authentication is required.")
		return
	}

	notifications, err := h.service.FindByRecipientID(userID)
	if err != nil {
		log.Printf("list current user's notifications: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, notifications)
}
