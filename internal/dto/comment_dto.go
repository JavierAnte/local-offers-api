package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	Body string `json:"body"`
}

type CommentResponse struct {
	ID uuid.UUID `json:"id"`

	OfferID uuid.UUID `json:"offerId"`

	Body string `json:"body"`

	CreatedAt time.Time    `json:"createdAt"`
	PostedBy  PostedByInfo `json:"postedBy"`
}
