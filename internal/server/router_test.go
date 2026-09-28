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
		Upload: &handlers.UploadHandler{}, Notification: &handlers.NotificationHandler{},
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
		Upload: &handlers.UploadHandler{}, Notification: &handlers.NotificationHandler{},
	})
	requests := []struct{ method, path string }{
		{http.MethodPost, "/api/v1/offers"},
		{http.MethodPost, "/api/v1/offers/id/comments"},
		{http.MethodPost, "/api/v1/offers/id/votes"},
		{http.MethodPost, "/api/v1/uploads"},
		{http.MethodGet, "/api/v1/me"},
		{http.MethodGet, "/api/v1/me/offers"},
		{http.MethodGet, "/api/v1/me/notifications"},
	}
	for _, request := range requests {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(request.method, request.path, nil))
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s status = %d", request.method, request.path, recorder.Code)
		}
		if !strings.Contains(recorder.Body.String(), `"code":"authentication_required"`) {
			t.Fatalf("%s %s body = %s", request.method, request.path, recorder.Body.String())
		}
	}
}
