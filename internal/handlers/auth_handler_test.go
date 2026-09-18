package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/services"
)

type fakeAuthService struct{ registerErr, loginErr error }

func (f *fakeAuthService) Register(dto.RegisterRequest) (*dto.AuthResponse, error) {
	return &dto.AuthResponse{}, f.registerErr
}
func (f *fakeAuthService) Login(dto.LoginRequest) (*dto.AuthResponse, error) {
	return &dto.AuthResponse{}, f.loginErr
}

func TestAuthHandlerStatuses(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		err    error
		status int
	}{
		{"register success", nil, 201}, {"register invalid", services.ErrInvalidInput, 400},
		{"register conflict", services.ErrEmailTaken, 409}, {"register internal", errors.New("db"), 500},
	} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"name":"User","email":"user@example.com","password":"password"}`))
			NewAuthHandler(&fakeAuthService{registerErr: tt.err}).Register(recorder, req)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.status)
			}
		})
	}

	for _, tt := range []struct {
		name   string
		err    error
		status int
	}{{"login success", nil, 200}, {"login invalid", services.ErrInvalidCredentials, 401}, {"login internal", errors.New("db"), 500}} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"user@example.com","password":"password"}`))
			NewAuthHandler(&fakeAuthService{loginErr: tt.err}).Login(recorder, req)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.status)
			}
		})
	}
}

func TestAuthHandlerRejectsUnknownAndTrailingJSON(t *testing.T) {
	t.Parallel()
	for _, body := range []string{
		`{"name":"User","email":"user@example.com","password":"password","admin":true}`,
		`{"name":"User","email":"user@example.com","password":"password"} {}`,
	} {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
		NewAuthHandler(&fakeAuthService{}).Register(recorder, req)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("status = %d: %s", recorder.Code, recorder.Body.String())
		}
		if !strings.Contains(recorder.Body.String(), `"code":"invalid_json"`) {
			t.Fatalf("body = %s", recorder.Body.String())
		}
	}
}
