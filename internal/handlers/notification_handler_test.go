package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/google/uuid"
)

type fakeNotificationService struct {
	notifications []dto.NotificationResponse
	err           error
	recipientID   uuid.UUID
}

func (f *fakeNotificationService) FindByRecipientID(recipientID uuid.UUID) ([]dto.NotificationResponse, error) {
	f.recipientID = recipientID
	return f.notifications, f.err
}

func TestNotificationHandlerList(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	token, err := auth.GenerateToken("secret", userID)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{name: "success", status: http.StatusOK},
		{name: "internal", err: errors.New("database unavailable"), status: http.StatusInternalServerError, code: "internal_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeNotificationService{notifications: []dto.NotificationResponse{}, err: tt.err}
			handler := NewNotificationHandler(service)
			req := httptest.NewRequest(http.MethodGet, "/me/notifications", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			auth.RequireAuth("secret")(http.HandlerFunc(handler.List)).ServeHTTP(recorder, req)

			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.status, recorder.Body.String())
			}
			if tt.code != "" && !strings.Contains(recorder.Body.String(), `"code":"`+tt.code+`"`) {
				t.Fatalf("body = %s", recorder.Body.String())
			}
			if tt.err == nil && service.recipientID != userID {
				t.Fatalf("recipient id = %s, want %s", service.recipientID, userID)
			}
		})
	}
}
