package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID uuid.UUID `json:"id"`

	Type string `json:"type"`

	OfferID       uuid.UUID    `json:"offerId"`
	OfferHeadline string       `json:"offerHeadline"`
	Actor         PostedByInfo `json:"actor"`

	CreatedAt time.Time `json:"createdAt"`
}
