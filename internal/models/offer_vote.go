package models

import (
	"time"

	"github.com/google/uuid"
)

type VoteType string

const (
	VoteTypeValidate   VoteType = "validate"
	VoteTypeInvalidate VoteType = "invalidate"
)

type OfferVote struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OfferID uuid.UUID `gorm:"type:uuid"`
	UserID  uuid.UUID `gorm:"type:uuid"`

	Type VoteType

	CreatedAt time.Time
}
