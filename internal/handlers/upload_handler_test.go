package handlers

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/services"
)

type fakeUploadService struct {
	filename string
	err      error
}

func (f *fakeUploadService) Save(io.Reader) (string, error) { return f.filename, f.err }

func multipartImageRequest(t *testing.T) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", "offer.jpg")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("image bytes"))
	_ = writer.Close()
	req := httptest.NewRequest(http.MethodPost, "http://api.test/uploads", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestUploadHandlerCreate(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name    string
		service *fakeUploadService
		request func(*testing.T) *http.Request
		status  int
	}{
		{"success", &fakeUploadService{filename: "stored.jpg"}, multipartImageRequest, 201},
		{"invalid image", &fakeUploadService{err: services.ErrInvalidInput}, multipartImageRequest, 400},
		{"storage error", &fakeUploadService{err: errors.New("disk")}, multipartImageRequest, 500},
		{"not multipart", &fakeUploadService{}, func(*testing.T) *http.Request {
			return httptest.NewRequest(http.MethodPost, "/uploads", strings.NewReader("bad"))
		}, 400},
	} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			NewUploadHandler(tt.service).Create(recorder, tt.request(t))
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, tt.status, recorder.Body.String())
			}
			if tt.name == "success" && !strings.Contains(recorder.Body.String(), "http://api.test/uploads/stored.jpg") {
				t.Fatalf("body = %s", recorder.Body.String())
			}
		})
	}
}
