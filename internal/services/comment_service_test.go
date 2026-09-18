package services

import (
	"errors"
	"strings"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
)

type fakeCommentRepository struct {
	exists  bool
	err     error
	created *models.Comment
}

func (f *fakeCommentRepository) Create(comment *models.Comment) error {
	f.created = comment
	return f.err
}
func (f *fakeCommentRepository) OfferExists(uuid.UUID) (bool, error) { return f.exists, f.err }
func (f *fakeCommentRepository) FindByOfferID(uuid.UUID) ([]dto.CommentResponse, error) {
	return []dto.CommentResponse{}, f.err
}

func TestCommentServiceCreate(t *testing.T) {
	t.Parallel()
	offerID := uuid.NewString()
	repo := &fakeCommentRepository{exists: true}
	if err := NewCommentService(repo).Create(dto.CreateCommentRequest{Body: "  still available  "}, offerID, uuid.New()); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if repo.created.Body != "still available" {
		t.Fatalf("body = %q", repo.created.Body)
	}
}

func TestCommentServiceRejectsInvalidOrMissingOffer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, id, body string
		exists         bool
		want           error
	}{
		{"invalid id", "bad", "hello", true, ErrInvalidInput},
		{"empty", uuid.NewString(), "  ", true, ErrInvalidInput},
		{"too long", uuid.NewString(), strings.Repeat("x", 1001), true, ErrInvalidInput},
		{"missing", uuid.NewString(), "hello", false, ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCommentService(&fakeCommentRepository{exists: tt.exists}).Create(dto.CreateCommentRequest{Body: tt.body}, tt.id, uuid.New())
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}
