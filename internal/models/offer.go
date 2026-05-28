package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Offer struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Headline    string
	Description *string

	BusinessName string
	Category     string

	ImageURL *string

	OfferType datatypes.JSON `gorm:"type:jsonb"`

	Location string `gorm:"type:geography(POINT,4326)"`

	ExpiresAt *time.Time

	ConfirmationsCount int
	InvalidationsCount int

	CreatedAt time.Time
}
