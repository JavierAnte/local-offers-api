package services

import (
	"errors"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeVoteRepository struct {
	confirmations, invalidations int
	err                          error
	voteType                     models.VoteType
}

func (f *fakeVoteRepository) Vote(_ uuid.UUID, _ uuid.UUID, voteType models.VoteType) (int, int, error) {
	f.voteType = voteType
	return f.confirmations, f.invalidations, f.err
}

func TestOfferVoteService(t *testing.T) {
	t.Parallel()
	repo := &fakeVoteRepository{confirmations: 2, invalidations: 1}
	got, err := NewOfferVoteService(repo).Vote(uuid.NewString(), uuid.New(), "validate")
	if err != nil {
		t.Fatalf("Vote() error = %v", err)
	}
	if got.ConfirmationsCount != 2 || got.InvalidationsCount != 1 {
		t.Fatalf("Vote() = %#v", got)
	}
	if repo.voteType != models.VoteTypeValidate {
		t.Fatalf("type = %q", repo.voteType)
	}
}

func TestOfferVoteServiceErrors(t *testing.T) {
	t.Parallel()
	service := NewOfferVoteService(&fakeVoteRepository{})
	if _, err := service.Vote("bad", uuid.New(), "validate"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("invalid id = %v", err)
	}
	if _, err := service.Vote(uuid.NewString(), uuid.New(), "maybe"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("invalid type = %v", err)
	}
	service = NewOfferVoteService(&fakeVoteRepository{err: gorm.ErrRecordNotFound})
	if _, err := service.Vote(uuid.NewString(), uuid.New(), "invalidate"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing = %v", err)
	}
}
