package repositories

import (
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// commentRow mirrors dto.CommentResponse's flat columns plus the joined
// poster fields, which are scanned separately and then folded into
// PostedBy since gorm's Raw+Scan doesn't populate nested structs from flat
// column aliases.
type commentRow struct {
	dto.CommentResponse
	PostedByID   *uuid.UUID
	PostedByName *string
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindByOfferID(offerID string) ([]dto.CommentResponse, error) {
	var rows []commentRow

	query := `
		SELECT
			c.id,
			c.offer_id,
			c.body,
			c.created_at,

			u.id AS posted_by_id,
			u.name AS posted_by_name
		FROM comments c
		LEFT JOIN users u ON u.id = c.user_id
		WHERE c.offer_id = ?
		ORDER BY c.created_at ASC;
	`

	err := r.db.Raw(query, offerID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	comments := make([]dto.CommentResponse, len(rows))
	for i := range rows {
		comments[i] = rows[i].CommentResponse
		comments[i].PostedBy = postedByFrom(rows[i].PostedByID, rows[i].PostedByName)
	}

	return comments, nil
}
