package services

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeOfferRepository struct {
	created   *models.Offer
	createErr error
	found     *dto.OfferResponse
	findErr   error
}

func (f *fakeOfferRepository) Create(offer *models.Offer) error {
	f.created = offer
	return f.createErr
}
func (f *fakeOfferRepository) FindNearby(float64, float64, int) ([]dto.OfferResponse, error) {
	return nil, f.findErr
}
func (f *fakeOfferRepository) FindByID(uuid.UUID) (*dto.OfferResponse, error) {
	return f.found, f.findErr
}

func validOfferRequest() dto.CreateOfferRequest {
	expires := time.Now().Add(time.Hour)
	return dto.CreateOfferRequest{
		Headline: "  Half-price lunch  ", BusinessName: "  Cafe Local  ", Category: "food",
		OfferType: json.RawMessage(`{"type":"percentage","percentage":50}`),
		ExpiresAt: &expires, Latitude: -31.2, Longitude: -64.4,
	}
}

func TestOfferServiceCreateValidatesAndNormalizes(t *testing.T) {
	t.Parallel()
	repo := &fakeOfferRepository{}
	err := NewOfferService(repo).Create(validOfferRequest(), uuid.New())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if repo.created == nil {
		t.Fatal("repository Create was not called")
	}
	if repo.created.Headline != "Half-price lunch" || repo.created.BusinessName != "Cafe Local" {
		t.Fatalf("strings were not normalized: %#v", repo.created)
	}
}

func TestOfferServiceRejectsInvalidOffers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(*dto.CreateOfferRequest)
	}{
		{"short headline", func(r *dto.CreateOfferRequest) { r.Headline = "bad" }},
		{"unknown category", func(r *dto.CreateOfferRequest) { r.Category = "cars" }},
		{"latitude", func(r *dto.CreateOfferRequest) { r.Latitude = 91 }},
		{"longitude", func(r *dto.CreateOfferRequest) { r.Longitude = -181 }},
		{"past expiry", func(r *dto.CreateOfferRequest) { v := time.Now().Add(-time.Minute); r.ExpiresAt = &v }},
		{"unknown offer field", func(r *dto.CreateOfferRequest) {
			r.OfferType = json.RawMessage(`{"type":"percentage","percentage":20,"extra":true}`)
		}},
		{"invalid percentage", func(r *dto.CreateOfferRequest) {
			r.OfferType = json.RawMessage(`{"type":"percentage","percentage":100}`)
		}},
		{"invalid price", func(r *dto.CreateOfferRequest) {
			r.OfferType = json.RawMessage(`{"type":"price","currentPrice":100,"previousPrice":50}`)
		}},
		{"empty label", func(r *dto.CreateOfferRequest) { r.OfferType = json.RawMessage(`{"type":"bundle","label":"  "}`) }},
		{"invalid image URL", func(r *dto.CreateOfferRequest) { v := "file:///tmp/image.jpg"; r.ImageURL = &v }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validOfferRequest()
			tt.mutate(&req)
			repo := &fakeOfferRepository{}
			if err := NewOfferService(repo).Create(req, uuid.New()); !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("Create() error = %v, want ErrInvalidInput", err)
			}
			if repo.created != nil {
				t.Fatal("invalid offer reached repository")
			}
		})
	}
}

func TestOfferServiceAcceptsSupportedOfferTypes(t *testing.T) {
	t.Parallel()
	values := []string{
		`{"type":"percentage","percentage":25}`,
		`{"type":"price","currentPrice":100,"previousPrice":150}`,
		`{"type":"bundle","label":"2x1"}`,
		`{"type":"text","label":"Cash only"}`,
	}
	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			req := validOfferRequest()
			req.OfferType = json.RawMessage(value)
			if err := NewOfferService(&fakeOfferRepository{}).Create(req, uuid.New()); err != nil {
				t.Fatalf("Create() error = %v", err)
			}
		})
	}
}

func TestOfferServiceFindByID(t *testing.T) {
	t.Parallel()
	service := NewOfferService(&fakeOfferRepository{findErr: gorm.ErrRecordNotFound})
	if _, err := service.FindByID("not-a-uuid"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("invalid id error = %v", err)
	}
	if _, err := service.FindByID(uuid.NewString()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing id error = %v", err)
	}
}
