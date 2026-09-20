package services

import (
	"errors"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeUserProfileRepository struct {
	user *models.User
	err  error
}

type fakeUserOffersRepository struct {
	offers []dto.OfferResponse
	err    error
}

func (f *fakeUserOffersRepository) FindByUserID(uuid.UUID) ([]dto.OfferResponse, error) {
	return f.offers, f.err
}

func (f *fakeUserProfileRepository) FindByID(uuid.UUID) (*models.User, error) {
	return f.user, f.err
}

func TestUserServiceMe(t *testing.T) {
	t.Parallel()
	user := &models.User{ID: uuid.New(), Name: "Local User", Email: "user@example.com"}
	got, err := NewUserService(&fakeUserProfileRepository{user: user}, &fakeUserOffersRepository{}).Me(user.ID)
	if err != nil {
		t.Fatalf("Me() error = %v", err)
	}
	if got.ID != user.ID || got.Name != user.Name || got.Email != user.Email {
		t.Fatalf("Me() = %#v", got)
	}
}

func TestUserServiceMeErrors(t *testing.T) {
	t.Parallel()
	if _, err := NewUserService(&fakeUserProfileRepository{err: gorm.ErrRecordNotFound}, &fakeUserOffersRepository{}).Me(uuid.New()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing user error = %v", err)
	}
	databaseErr := errors.New("database unavailable")
	if _, err := NewUserService(&fakeUserProfileRepository{err: databaseErr}, &fakeUserOffersRepository{}).Me(uuid.New()); !errors.Is(err, databaseErr) {
		t.Fatalf("database error = %v", err)
	}
}

func TestUserServiceMyOffers(t *testing.T) {
	t.Parallel()
	offers := []dto.OfferResponse{{ID: uuid.New(), Headline: "My offer"}}
	got, err := NewUserService(&fakeUserProfileRepository{}, &fakeUserOffersRepository{offers: offers}).MyOffers(uuid.New())
	if err != nil || len(got) != 1 || got[0].ID != offers[0].ID {
		t.Fatalf("MyOffers() = %#v, %v", got, err)
	}
}
