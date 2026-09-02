package repositories

import (
	"math/rand"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfferRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

// offerRow mirrors dto.OfferResponse's flat columns plus the joined poster
// fields, which are scanned separately and then folded into PostedBy since
// gorm's Raw+Scan doesn't populate nested structs from flat column aliases.
type offerRow struct {
	dto.OfferResponse
	PostedByID   *uuid.UUID
	PostedByName *string
}

func (r *OfferRepository) Create(offer *models.Offer) error {
	return r.db.Create(offer).Error
}

func (r *OfferRepository) FindNearby(
	latitude float64,
	longitude float64,
	radiusMeters int,
) ([]dto.OfferResponse, error) {

	var rows []offerRow

	query := `
		SELECT
			o.id,
			o.headline,
			o.description,
			o.business_name,
			o.category,
			o.image_url,
			o.offer_type,
			o.expires_at,

			ST_Y(o.location::geometry) AS latitude,
			ST_X(o.location::geometry) AS longitude,

			ST_Distance(
				o.location,
				ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography
			) AS distance_meters,

			o.confirmations_count,
			o.invalidations_count,
			o.created_at,

			u.id AS posted_by_id,
			u.name AS posted_by_name
		FROM offers o
		LEFT JOIN users u ON u.id = o.user_id
		WHERE ST_DWithin(
			o.location,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
		ORDER BY distance_meters ASC
		LIMIT 50;
	`

	err := r.db.Raw(
		query,
		longitude,
		latitude,
		longitude,
		latitude,
		radiusMeters,
	).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	offers := make([]dto.OfferResponse, len(rows))
	for i := range rows {
		offers[i] = rows[i].OfferResponse
		offers[i].IsVerifiedBusiness = false
		offers[i].CommentsCount = rand.Intn(11) // 0-10 — comments aren't modeled yet
		offers[i].PostedBy = postedByFrom(rows[i].PostedByID, rows[i].PostedByName)
	}

	return offers, nil
}

func (r *OfferRepository) FindByID(id string) (*dto.OfferResponse, error) {
	var row offerRow

	query := `
		SELECT
			o.id,
			o.headline,
			o.description,
			o.business_name,
			o.category,
			o.image_url,
			o.offer_type,
			o.expires_at,

			ST_Y(o.location::geometry) AS latitude,
			ST_X(o.location::geometry) AS longitude,

			o.confirmations_count,
			o.invalidations_count,
			o.created_at,

			u.id AS posted_by_id,
			u.name AS posted_by_name
		FROM offers o
		LEFT JOIN users u ON u.id = o.user_id
		WHERE o.id = ?
		LIMIT 1;
	`

	err := r.db.Raw(query, id).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	offer := row.OfferResponse
	offer.IsVerifiedBusiness = false
	offer.CommentsCount = rand.Intn(11) // 0-10 — comments aren't modeled yet
	offer.PostedBy = postedByFrom(row.PostedByID, row.PostedByName)

	return &offer, nil
}

func postedByFrom(id *uuid.UUID, name *string) dto.PostedByInfo {
	if id == nil {
		return dto.PostedByInfo{Name: "Usuario LocalOffers"}
	}

	posterName := "Usuario LocalOffers"
	if name != nil {
		posterName = *name
	}

	return dto.PostedByInfo{ID: id.String(), Name: posterName}
}
