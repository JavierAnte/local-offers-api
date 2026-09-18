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

type fakeOfferService struct {
	findResult          *dto.OfferResponse
	err                 error
	latitude, longitude float64
	radius              int
}

func TestOfferHandlerCreate(t *testing.T) {
	t.Parallel()
	token, _ := auth.GenerateToken("secret", uuid.New())
	for _, tt := range []struct {
		name, body string
		err        error
		status     int
	}{
		{"success", `{}`, nil, 201}, {"bad json", `{"extra":true}`, nil, 400},
		{"invalid", `{}`, services.ErrInvalidInput, 400}, {"internal", `{}`, errors.New("db"), 500},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := NewOfferHandler(&fakeOfferService{err: tt.err})
			r := chi.NewRouter()
			r.With(auth.RequireAuth("secret")).Post("/offers", h.Create)
			req := httptest.NewRequest(http.MethodPost, "/offers", strings.NewReader(tt.body))
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.status, recorder.Body.String())
			}
		})
	}
}

func (f *fakeOfferService) Create(dto.CreateOfferRequest, uuid.UUID) error { return f.err }
func (f *fakeOfferService) FindNearby(lat, lng float64, radius int) ([]dto.OfferResponse, error) {
	f.latitude, f.longitude, f.radius = lat, lng, radius
	return []dto.OfferResponse{}, f.err
}
func (f *fakeOfferService) FindByID(string) (*dto.OfferResponse, error) { return f.findResult, f.err }

func offerRouter(service offerService) http.Handler {
	r := chi.NewRouter()
	h := NewOfferHandler(service)
	r.Get("/offers/nearby", h.FindNearby)
	r.Get("/offers/{id}", h.FindByID)
	return r
}

func TestOfferHandlerFindNearby(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, query            string
		wantStatus, wantRadius int
	}{
		{"default radius", "?lat=-31.2&lng=-64.4", 200, 5000},
		{"custom radius", "?lat=-31.2&lng=-64.4&radius=1000", 200, 1000},
		{"clamped radius", "?lat=-31.2&lng=-64.4&radius=99999", 200, 20000},
		{"invalid radius", "?lat=-31.2&lng=-64.4&radius=0", 400, 0},
		{"invalid coordinate", "?lat=nope&lng=-64.4", 400, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := &fakeOfferService{}
			recorder := httptest.NewRecorder()
			offerRouter(service).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/offers/nearby"+tt.query, nil))
			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if tt.wantRadius != 0 && service.radius != tt.wantRadius {
				t.Fatalf("radius = %d, want %d", service.radius, tt.wantRadius)
			}
		})
	}
}

func TestOfferHandlerFindByIDErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{"invalid", services.ErrInvalidInput, 400, "invalid_offer_id"},
		{"missing", services.ErrNotFound, 404, "offer_not_found"},
		{"internal", errors.New("database unavailable"), 500, "internal_error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			offerRouter(&fakeOfferService{err: tt.err}).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/offers/"+uuid.NewString(), nil))
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.status)
			}
			if !strings.Contains(recorder.Body.String(), `"code":"`+tt.code+`"`) {
				t.Fatalf("body = %s", recorder.Body.String())
			}
		})
	}
}
