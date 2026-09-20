package repositories

import (
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
