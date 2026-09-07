package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OfferID uuid.UUID `gorm:"type:uuid"`
	UserID  uuid.UUID `gorm:"type:uuid"`

	Body string

	CreatedAt time.Time
}
