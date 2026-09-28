package services

import (
	"errors"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/google/uuid"
)

type fakeNotificationRepository struct {
	notifications []dto.NotificationResponse
	err           error
	recipientID   uuid.UUID
}

func (f *fakeNotificationRepository) FindByRecipientID(recipientID uuid.UUID) ([]dto.NotificationResponse, error) {
	f.recipientID = recipientID
	return f.notifications, f.err
}

func TestNotificationServiceFindByRecipientID(t *testing.T) {
	t.Parallel()
	recipientID := uuid.New()
	want := []dto.NotificationResponse{{ID: uuid.New(), Type: "comment_received"}}
	repo := &fakeNotificationRepository{notifications: want}

	got, err := NewNotificationService(repo).FindByRecipientID(recipientID)
	if err != nil {
		t.Fatalf("FindByRecipientID() error = %v", err)
	}
	if repo.recipientID != recipientID || len(got) != 1 || got[0].ID != want[0].ID {
		t.Fatalf("FindByRecipientID() = %#v", got)
	}

	wantErr := errors.New("database unavailable")
	if _, err := NewNotificationService(&fakeNotificationRepository{err: wantErr}).FindByRecipientID(recipientID); !errors.Is(err, wantErr) {
		t.Fatalf("FindByRecipientID() error = %v, want %v", err, wantErr)
	}
}
