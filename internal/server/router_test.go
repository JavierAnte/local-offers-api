package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/handlers"
)

func TestRouterHealth(t *testing.T) {
	t.Parallel()
	router := NewRouter("secret", t.TempDir(), Handlers{
		Auth: &handlers.AuthHandler{}, Offer: &handlers.OfferHandler{},
		Comment: &handlers.CommentHandler{}, OfferVote: &handlers.OfferVoteHandler{},
		Upload: &handlers.UploadHandler{},
	})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))
	if recorder.Code != http.StatusOK || recorder.Body.String() != "OK" {
		t.Fatalf("health response = %d %q", recorder.Code, recorder.Body.String())
	}
}

func TestRouterProtectsWriteRoutes(t *testing.T) {
	t.Parallel()
	router := NewRouter("secret", t.TempDir(), Handlers{
		Auth: &handlers.AuthHandler{}, Offer: &handlers.OfferHandler{},
		Comment: &handlers.CommentHandler{}, OfferVote: &handlers.OfferVoteHandler{},
		Upload: &handlers.UploadHandler{},
	})
	paths := []string{"/api/v1/offers", "/api/v1/offers/id/comments", "/api/v1/offers/id/votes", "/api/v1/uploads"}
	for _, path := range paths {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, path, nil))
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("POST %s status = %d", path, recorder.Code)
		}
		if !strings.Contains(recorder.Body.String(), `"code":"authentication_required"`) {
			t.Fatalf("POST %s body = %s", path, recorder.Body.String())
		}
	}
}
