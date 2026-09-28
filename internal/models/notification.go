package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeCommentReceived  NotificationType = "comment_received"
	NotificationTypeOfferValidated   NotificationType = "offer_validated"
	NotificationTypeOfferInvalidated NotificationType = "offer_invalidated"
)

type Notification struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	RecipientUserID uuid.UUID `gorm:"type:uuid"`
	ActorUserID     uuid.UUID `gorm:"type:uuid"`
	OfferID         uuid.UUID `gorm:"type:uuid"`

	Type NotificationType

	CreatedAt time.Time
}
