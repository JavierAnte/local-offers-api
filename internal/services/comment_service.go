package services

import (
	"errors"
	"strings"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
	repo commentRepository
}

type commentRepository interface {
	Create(comment *models.Comment) error
	OfferExists(offerID uuid.UUID) (bool, error)
	FindByOfferID(offerID uuid.UUID) ([]dto.CommentResponse, error)
}

func NewCommentService(repo commentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(req dto.CreateCommentRequest, offerID string, userID uuid.UUID) error {
	body := strings.TrimSpace(req.Body)
	if body == "" || len(body) > 1000 {
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

		Body: body,

		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&comment); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *CommentService) FindByOfferID(offerID string) ([]dto.CommentResponse, error) {
	parsedOfferID, err := uuid.Parse(offerID)
	if err != nil {
		return nil, ErrInvalidInput
	}
	exists, err := s.repo.OfferExists(parsedOfferID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound
	}
	return s.repo.FindByOfferID(parsedOfferID)
}
