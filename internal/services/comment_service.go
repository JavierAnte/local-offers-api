package services

import (
	"strings"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"

	"github.com/google/uuid"
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
	exists, err := s.repo.OfferExists(parsedOfferID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	comment := models.Comment{
		ID: uuid.New(),

		OfferID: parsedOfferID,
		UserID:  userID,

		Body: body,

		CreatedAt: time.Now(),
	}

	return s.repo.Create(&comment)
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
