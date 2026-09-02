package services

import (
	"errors"
	"strings"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/auth"
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/JavierAnte/local-offers-api/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo      *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(repo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if name == "" || email == "" || len(req.Password) < 8 {
		return nil, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	return s.buildAuthResponse(user)
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	return s.buildAuthResponse(user)
}

func (s *AuthService) buildAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	token, err := auth.GenerateToken(s.jwtSecret, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
