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

type fakeCommentService struct{ err error }

func (f *fakeCommentService) Create(dto.CreateCommentRequest, string, uuid.UUID) error { return f.err }
func (f *fakeCommentService) FindByOfferID(string) ([]dto.CommentResponse, error) {
	return []dto.CommentResponse{}, f.err
}

func commentRouter(t *testing.T, service commentService) http.Handler {
	t.Helper()
	h := NewCommentHandler(service)
	r := chi.NewRouter()
	r.Get("/offers/{id}/comments", h.List)
	r.With(auth.RequireAuth("secret")).Post("/offers/{id}/comments", h.Create)
	return r
}

func TestCommentHandlerStatuses(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		err    error
		status int
	}{
		{"success", nil, 200}, {"invalid", services.ErrInvalidInput, 400},
		{"missing", services.ErrNotFound, 404}, {"internal", errors.New("db"), 500},
	} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			commentRouter(t, &fakeCommentService{err: tt.err}).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/offers/"+uuid.NewString()+"/comments", nil))
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.status)
			}
		})
	}
}

func TestCommentHandlerCreate(t *testing.T) {
	t.Parallel()
	token, _ := auth.GenerateToken("secret", uuid.New())
	for _, tt := range []struct {
		name, body string
		err        error
		status     int
	}{
		{"success", `{"body":"available"}`, nil, 201},
		{"bad json", `{"body":"ok","extra":true}`, nil, 400},
		{"invalid", `{"body":""}`, services.ErrInvalidInput, 400},
		{"missing", `{"body":"ok"}`, services.ErrNotFound, 404},
		{"internal", `{"body":"ok"}`, errors.New("db"), 500},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/offers/"+uuid.NewString()+"/comments", strings.NewReader(tt.body))
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			commentRouter(t, &fakeCommentService{err: tt.err}).ServeHTTP(recorder, req)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.status, recorder.Body.String())
			}
		})
	}
}
