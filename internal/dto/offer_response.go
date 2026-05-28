package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PostedByInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OfferResponse struct {
	ID uuid.UUID `json:"id"`

	Headline    string  `json:"headline"`
	Description *string `json:"description"`

	BusinessName       string `json:"businessName"`
	IsVerifiedBusiness bool   `json:"isVerifiedBusiness"`
	Category           string `json:"category"`

	ImageURL *string `json:"imageUrl"`

	OfferType datatypes.JSON `json:"offerType"`

	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	DistanceMeters float64 `json:"distanceMeters,omitempty"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	ConfirmationsCount int `json:"confirmationsCount"`
	InvalidationsCount int `json:"invalidationsCount"`
	CommentsCount      int `json:"commentsCount"`

	CreatedAt time.Time    `json:"createdAt"`
	PostedBy  PostedByInfo `json:"postedBy" gorm:"-"`
}
