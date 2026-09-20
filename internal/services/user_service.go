package services

import (
	"errors"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo  userProfileRepository
	offerRepo userOffersRepository
}

type userProfileRepository interface {
	FindByID(id uuid.UUID) (*models.User, error)
}

type userOffersRepository interface {
	FindByUserID(userID uuid.UUID) ([]dto.OfferResponse, error)
}

func NewUserService(userRepo userProfileRepository, offerRepo userOffersRepository) *UserService {
	return &UserService{userRepo: userRepo, offerRepo: offerRepo}
}

func (s *UserService) Me(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}

func (s *UserService) MyOffers(userID uuid.UUID) ([]dto.OfferResponse, error) {
	return s.offerRepo.FindByUserID(userID)
}
