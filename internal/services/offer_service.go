package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OfferService struct {
	repo offerRepository
}

type offerRepository interface {
	Create(offer *models.Offer) error
	FindNearby(latitude float64, longitude float64, radiusMeters int) ([]dto.OfferResponse, error)
	FindByID(id uuid.UUID) (*dto.OfferResponse, error)
}

func NewOfferService(repo offerRepository) *OfferService {
	return &OfferService{repo: repo}
}

func (s *OfferService) Create(req dto.CreateOfferRequest, userID uuid.UUID) error {
	if err := validateAndNormalizeOffer(&req); err != nil {
		return err
	}

	location := fmt.Sprintf(
		"SRID=4326;POINT(%f %f)",
		req.Longitude,
		req.Latitude,
	)

	offer := models.Offer{
		ID: uuid.New(),

		Headline:    req.Headline,
		Description: req.Description,

		BusinessName: req.BusinessName,
		Category:     req.Category,

		ImageURL: req.ImageURL,

		OfferType: datatypes.JSON(req.OfferType),

		Location: location,

		UserID: &userID,

		ExpiresAt: req.ExpiresAt,

		CreatedAt: time.Now(),
	}

	return s.repo.Create(&offer)
}

func (s *OfferService) FindNearby(
	latitude float64,
	longitude float64,
	radiusMeters int,
) ([]dto.OfferResponse, error) {
	if !validCoordinates(latitude, longitude) || radiusMeters <= 0 {
		return nil, ErrInvalidInput
	}

	return s.repo.FindNearby(
		latitude,
		longitude,
		radiusMeters,
	)
}

func (s *OfferService) FindByID(id string) (*dto.OfferResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidInput
	}

	offer, err := s.repo.FindByID(parsedID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return offer, err
}

var validCategories = map[string]struct{}{
	"food": {}, "grocery": {}, "electronics": {}, "fashion": {},
	"beauty": {}, "sports": {}, "home": {}, "other": {},
}

type offerTypePayload struct {
	Type          string   `json:"type"`
	Percentage    *float64 `json:"percentage,omitempty"`
	CurrentPrice  *float64 `json:"currentPrice,omitempty"`
	PreviousPrice *float64 `json:"previousPrice,omitempty"`
	Label         *string  `json:"label,omitempty"`
}

func validateAndNormalizeOffer(req *dto.CreateOfferRequest) error {
	req.Headline = strings.TrimSpace(req.Headline)
	req.BusinessName = strings.TrimSpace(req.BusinessName)
	req.Category = strings.TrimSpace(req.Category)

	if len(req.Headline) < 5 || len(req.Headline) > 100 || len(req.BusinessName) < 2 || len(req.BusinessName) > 80 {
		return ErrInvalidInput
	}
	if _, ok := validCategories[req.Category]; !ok {
		return ErrInvalidInput
	}
	if !validCoordinates(req.Latitude, req.Longitude) {
		return ErrInvalidInput
	}
	if req.Description != nil {
		value := strings.TrimSpace(*req.Description)
		if len(value) > 1000 {
			return ErrInvalidInput
		}
		req.Description = &value
	}
	if req.ImageURL != nil {
		value := strings.TrimSpace(*req.ImageURL)
		parsed, err := url.ParseRequestURI(value)
		if err != nil || len(value) > 2048 || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return ErrInvalidInput
		}
		req.ImageURL = &value
	}
	if req.ExpiresAt != nil && !req.ExpiresAt.After(time.Now()) {
		return ErrInvalidInput
	}

	normalized, err := validateOfferType(req.OfferType)
	if err != nil {
		return err
	}
	req.OfferType = normalized
	return nil
}

func validCoordinates(latitude float64, longitude float64) bool {
	return !math.IsNaN(latitude) && !math.IsInf(latitude, 0) && latitude >= -90 && latitude <= 90 &&
		!math.IsNaN(longitude) && !math.IsInf(longitude, 0) && longitude >= -180 && longitude <= 180
}

func validateOfferType(raw json.RawMessage) (json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var value offerTypePayload
	if err := decoder.Decode(&value); err != nil {
		return nil, ErrInvalidInput
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, ErrInvalidInput
	}

	valid := false
	switch value.Type {
	case "percentage":
		valid = value.Percentage != nil && *value.Percentage > 0 && *value.Percentage < 100 && value.CurrentPrice == nil && value.PreviousPrice == nil && value.Label == nil
	case "price":
		valid = value.CurrentPrice != nil && *value.CurrentPrice > 0 && value.Percentage == nil && value.Label == nil &&
			(value.PreviousPrice == nil || (*value.PreviousPrice > 0 && *value.PreviousPrice > *value.CurrentPrice))
	case "bundle", "text":
		if value.Label != nil {
			label := strings.TrimSpace(*value.Label)
			value.Label = &label
			valid = label != "" && len(label) <= 80 && value.Percentage == nil && value.CurrentPrice == nil && value.PreviousPrice == nil
		}
	}
	if !valid {
		return nil, ErrInvalidInput
	}

	normalized, err := json.Marshal(value)
	if err != nil {
		return nil, ErrInvalidInput
	}
	return normalized, nil
}
