package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type fakeUserService struct {
	user   *dto.UserResponse
	offers []dto.OfferResponse
	err    error
}

func (f *fakeUserService) MyOffers(uuid.UUID) ([]dto.OfferResponse, error) {
	return f.offers, f.err
}

func (f *fakeUserService) Me(uuid.UUID) (*dto.UserResponse, error) {
	return f.user, f.err
}

func TestUserHandlerMyOffers(t *testing.T) {
	t.Parallel()
	token, err := auth.GenerateToken("secret", uuid.New())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		err    error
		status int
	}{
		{name: "success", status: 200},
		{name: "internal", err: errors.New("database unavailable"), status: 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewUserHandler(&fakeUserService{offers: []dto.OfferResponse{}, err: tt.err})
			r := chi.NewRouter()
			r.With(auth.RequireAuth("secret")).Get("/me/offers", h.MyOffers)
			req := httptest.NewRequest(http.MethodGet, "/me/offers", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.status, recorder.Body.String())
			}
		})
	}
}

func TestUserHandlerMe(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	token, err := auth.GenerateToken("secret", userID)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		user   *dto.UserResponse
		err    error
		status int
		code   string
	}{
		{name: "success", user: &dto.UserResponse{ID: userID, Name: "User", Email: "user@example.com"}, status: 200},
		{name: "missing", err: services.ErrNotFound, status: 404, code: "user_not_found"},
		{name: "internal", err: errors.New("database unavailable"), status: 500, code: "internal_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewUserHandler(&fakeUserService{user: tt.user, err: tt.err})
			r := chi.NewRouter()
			r.With(auth.RequireAuth("secret")).Get("/me", h.Me)
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.status, recorder.Body.String())
			}
			if tt.code != "" && !strings.Contains(recorder.Body.String(), `"code":"`+tt.code+`"`) {
				t.Fatalf("body = %s", recorder.Body.String())
			}
		})
	}
}
