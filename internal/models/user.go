package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Name         string
	Email        string
	PasswordHash string

	CreatedAt time.Time
}
