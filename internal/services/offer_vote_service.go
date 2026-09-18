package services

import (
	"errors"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfferVoteService struct {
	repo offerVoteRepository
}

type offerVoteRepository interface {
	Vote(offerID uuid.UUID, userID uuid.UUID, voteType models.VoteType) (int, int, error)
}

func NewOfferVoteService(repo offerVoteRepository) *OfferVoteService {
	return &OfferVoteService{repo: repo}
}

func (s *OfferVoteService) Vote(offerID string, userID uuid.UUID, voteType string) (dto.VoteResponse, error) {
	parsedOfferID, err := uuid.Parse(offerID)
	if err != nil {
		return dto.VoteResponse{}, ErrInvalidInput
	}

	parsedType := models.VoteType(voteType)
	if parsedType != models.VoteTypeValidate && parsedType != models.VoteTypeInvalidate {
		return dto.VoteResponse{}, ErrInvalidInput
	}

	confirmations, invalidations, err := s.repo.Vote(parsedOfferID, userID, parsedType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.VoteResponse{}, ErrNotFound
		}
		return dto.VoteResponse{}, err
	}

	return dto.VoteResponse{
		ConfirmationsCount: confirmations,
		InvalidationsCount: invalidations,
	}, nil
}
