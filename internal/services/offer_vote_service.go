package services

import (
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/JavierAnte/local-offers-api/internal/repositories"

	"github.com/google/uuid"
)

type OfferVoteService struct {
	repo *repositories.OfferVoteRepository
}

func NewOfferVoteService(repo *repositories.OfferVoteRepository) *OfferVoteService {
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
		return dto.VoteResponse{}, err
	}

	return dto.VoteResponse{
		ConfirmationsCount: confirmations,
		InvalidationsCount: invalidations,
	}, nil
}
