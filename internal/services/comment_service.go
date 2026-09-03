package services

import (
	"strings"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/JavierAnte/local-offers-api/internal/repositories"

	"github.com/google/uuid"
)

type CommentService struct {
	repo *repositories.CommentRepository
}

func NewCommentService(repo *repositories.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(req dto.CreateCommentRequest, offerID string, userID uuid.UUID) error {
	if strings.TrimSpace(req.Body) == "" {
		return ErrInvalidInput
	}

	parsedOfferID, err := uuid.Parse(offerID)
	if err != nil {
		return ErrInvalidInput
	}

	comment := models.Comment{
		ID: uuid.New(),

		OfferID: parsedOfferID,
		UserID:  userID,

		Body: req.Body,

		CreatedAt: time.Now(),
	}

	return s.repo.Create(&comment)
}

func (s *CommentService) FindByOfferID(offerID string) ([]dto.CommentResponse, error) {
	return s.repo.FindByOfferID(offerID)
}
