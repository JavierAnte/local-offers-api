package repositories

import (
	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OfferVoteRepository struct {
	db *gorm.DB
}

func NewOfferVoteRepository(db *gorm.DB) *OfferVoteRepository {
	return &OfferVoteRepository{db: db}
}

// Vote creates or updates a user's vote on an offer, then recalculates and
// persists the offer's confirmations/invalidations counters from the
// resulting offer_votes rows, all inside one transaction so the counters
// never drift from the underlying votes.
func (r *OfferVoteRepository) Vote(
	offerID uuid.UUID,
	userID uuid.UUID,
	voteType models.VoteType,
) (confirmationsCount int, invalidationsCount int, err error) {

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var offer models.Offer
		// Serialize votes for the same offer so each transaction counts all
		// previously committed votes before refreshing the denormalized totals.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id", "user_id").First(&offer, "id = ?", offerID).Error; err != nil {
			return err
		}

		var previousVote models.OfferVote
		previousVoteResult := tx.Where("offer_id = ? AND user_id = ?", offerID, userID).Limit(1).Find(&previousVote)
		if previousVoteResult.Error != nil {
			return previousVoteResult.Error
		}
		voteChanged := previousVoteResult.RowsAffected == 0 || previousVote.Type != voteType

		vote := models.OfferVote{
			ID: uuid.New(),

			OfferID: offerID,
			UserID:  userID,

			Type: voteType,
		}

		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "offer_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"type", "created_at"}),
		}).Create(&vote).Error
		if err != nil {
			return err
		}

		var confirmations, invalidations int64

		err = tx.Model(&models.OfferVote{}).
			Where("offer_id = ? AND type = ?", offerID, models.VoteTypeValidate).
			Count(&confirmations).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.OfferVote{}).
			Where("offer_id = ? AND type = ?", offerID, models.VoteTypeInvalidate).
			Count(&invalidations).Error
		if err != nil {
			return err
		}

		confirmationsCount = int(confirmations)
		invalidationsCount = int(invalidations)

		if err := tx.Model(&models.Offer{}).Where("id = ?", offerID).Updates(map[string]interface{}{
			"confirmations_count": confirmationsCount,
			"invalidations_count": invalidationsCount,
		}).Error; err != nil {
			return err
		}

		if !voteChanged || offer.UserID == nil || *offer.UserID == userID {
			return nil
		}

		notificationType := models.NotificationTypeOfferValidated
		if voteType == models.VoteTypeInvalidate {
			notificationType = models.NotificationTypeOfferInvalidated
		}
		notification := models.Notification{
			ID:              uuid.New(),
			RecipientUserID: *offer.UserID,
			ActorUserID:     userID,
			OfferID:         offerID,
			Type:            notificationType,
		}
		return tx.Create(&notification).Error
	})

	return confirmationsCount, invalidationsCount, err
}
