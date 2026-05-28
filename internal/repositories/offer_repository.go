package repositories

import (
	"math/rand"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"gorm.io/gorm"
)

type OfferRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

func (r *OfferRepository) Create(offer *models.Offer) error {
	return r.db.Create(offer).Error
}

func (r *OfferRepository) FindNearby(
	latitude float64,
	longitude float64,
	radiusMeters int,
) ([]dto.OfferResponse, error) {

	var offers []dto.OfferResponse

	query := `
		SELECT
			id,
			headline,
			description,
			business_name,
			category,
			image_url,
			offer_type,
			expires_at,

			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,

			ST_Distance(
				location,
				ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography
			) AS distance_meters,

			confirmations_count,
			invalidations_count,
			created_at
		FROM offers
		WHERE ST_DWithin(
			location,
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
	).Scan(&offers).Error

	if err != nil {
		return nil, err
	}

	// Populate mocked fields
	for i := range offers {
		offers[i].IsVerifiedBusiness = false
		offers[i].CommentsCount = rand.Intn(11) // 0-10
		offers[i].PostedBy = dto.PostedByInfo{
			ID:   1,
			Name: "Juan Perez",
		}
	}

	return offers, nil
}

func (r *OfferRepository) FindByID(id string) (*dto.OfferResponse, error) {
	var offer dto.OfferResponse

	query := `
		SELECT
			id,
			headline,
			description,
			business_name,
			category,
			image_url,
			offer_type,
			expires_at,

			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,

			confirmations_count,
			invalidations_count,
			created_at
		FROM offers
		WHERE id = ?
		LIMIT 1;
	`

	err := r.db.Raw(query, id).Scan(&offer).Error
	if err != nil {
		return nil, err
	}

	// Populate mocked fields
	offer.IsVerifiedBusiness = false
	offer.CommentsCount = rand.Intn(11) // 0-10
	offer.PostedBy = dto.PostedByInfo{
		ID:   1,
		Name: "Juan Perez",
	}

	return &offer, nil
}
