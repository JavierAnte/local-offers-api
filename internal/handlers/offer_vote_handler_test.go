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

type fakeOfferVoteService struct{ err error }

func (f *fakeOfferVoteService) Vote(string, uuid.UUID, string) (dto.VoteResponse, error) {
	return dto.VoteResponse{ConfirmationsCount: 1}, f.err
}

func TestOfferVoteHandlerCreate(t *testing.T) {
	t.Parallel()
	token, _ := auth.GenerateToken("secret", uuid.New())
	for _, tt := range []struct {
		name, body string
		err        error
		status     int
	}{
		{"success", `{"type":"validate"}`, nil, 200}, {"bad json", `{`, nil, 400},
		{"invalid", `{"type":"maybe"}`, services.ErrInvalidInput, 400},
		{"missing", `{"type":"validate"}`, services.ErrNotFound, 404},
		{"internal", `{"type":"validate"}`, errors.New("db"), 500},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := NewOfferVoteHandler(&fakeOfferVoteService{err: tt.err})
			r := chi.NewRouter()
			r.With(auth.RequireAuth("secret")).Post("/offers/{id}/votes", h.Create)
			req := httptest.NewRequest(http.MethodPost, "/offers/"+uuid.NewString()+"/votes", strings.NewReader(tt.body))
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.status, recorder.Body.String())
			}
		})
	}
}
