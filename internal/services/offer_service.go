package services

import (
	"fmt"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/JavierAnte/local-offers-api/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OfferService struct {
	repo *repositories.OfferRepository
}

func NewOfferService(repo *repositories.OfferRepository) *OfferService {
	return &OfferService{repo: repo}
}

func (s *OfferService) Create(req dto.CreateOfferRequest, userID uuid.UUID) error {
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

	return s.repo.FindNearby(
		latitude,
		longitude,
		radiusMeters,
	)
}

func (s *OfferService) FindByID(id string) (*dto.OfferResponse, error) {
	return s.repo.FindByID(id)
}
