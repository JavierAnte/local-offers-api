package dto

import "encoding/json"

type CreateOfferRequest struct {
	Headline    string  `json:"headline"`
	Description *string `json:"description"`

	BusinessName string `json:"businessName"`
	Category     string `json:"category"`

	ImageURL *string `json:"imageUrl"`

	OfferType json.RawMessage `json:"offerType"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
