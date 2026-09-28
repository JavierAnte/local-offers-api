package dto

type NearbyOffersQuery struct {
	Latitude     float64
	Longitude    float64
	RadiusMeters int
	Category     string
	Search       string
}
