package repositories

import (
	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

type notificationRow struct {
	dto.NotificationResponse
	ActorID   *uuid.UUID
	ActorName *string
}

func (r *NotificationRepository) FindByRecipientID(recipientID uuid.UUID) ([]dto.NotificationResponse, error) {
	var rows []notificationRow

	query := `
		SELECT
			n.id,
			n.type,
			n.offer_id,
			o.headline AS offer_headline,
			n.created_at,
			u.id AS actor_id,
			u.name AS actor_name
		FROM notifications n
		JOIN offers o ON o.id = n.offer_id
		JOIN users u ON u.id = n.actor_user_id
		WHERE n.recipient_user_id = ?
		ORDER BY n.created_at DESC, n.id DESC
		LIMIT 50;
	`

	if err := r.db.Raw(query, recipientID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	notifications := make([]dto.NotificationResponse, len(rows))
	for i := range rows {
		notifications[i] = rows[i].NotificationResponse
		notifications[i].Actor = postedByFrom(rows[i].ActorID, rows[i].ActorName)
	}

	return notifications, nil
}
