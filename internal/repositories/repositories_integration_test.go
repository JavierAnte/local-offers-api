//go:build integration

package repositories

import (
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func integrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("LOCAL_OFFERS_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("LOCAL_OFFERS_TEST_DATABASE_URL is not set")
	}
	parsed, err := url.Parse(dsn)
	if err != nil || !strings.Contains(strings.ToLower(strings.TrimPrefix(parsed.Path, "/")), "test") {
		t.Fatalf("refusing to reset database whose name does not contain 'test': %q", dsn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("connect to integration database: %v", err)
	}
	if err := db.Exec("DROP TABLE IF EXISTS offer_votes, comments, offers, users CASCADE").Error; err != nil {
		t.Fatalf("reset integration database: %v", err)
	}

	files, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.up.sql"))
	if err != nil {
		t.Fatalf("find migrations: %v", err)
	}
	sort.Strings(files)
	for _, path := range files {
		migration, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read migration %s: %v", path, err)
		}
		if err := db.Exec(string(migration)).Error; err != nil {
			t.Fatalf("apply migration %s: %v", path, err)
		}
	}
	return db
}

func TestRepositoriesWithPostGIS(t *testing.T) {
	db := integrationDB(t)
	userID := uuid.New()
	offerNearID := uuid.New()
	offerFarID := uuid.New()

	if err := db.Exec(`INSERT INTO users (id, name, email, password_hash) VALUES (?, 'Tester', 'tester@example.com', 'hash')`, userID).Error; err != nil {
		t.Fatalf("insert user: %v", err)
	}
	insertOffer := `INSERT INTO offers (id, headline, business_name, category, offer_type, location, user_id)
		VALUES (?, ?, 'Local', 'food', '{"type":"percentage","percentage":20}', ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)`
	if err := db.Exec(insertOffer, offerNearID, "Near", -64.4, -31.2, userID).Error; err != nil {
		t.Fatalf("insert near offer: %v", err)
	}
	if err := db.Exec(insertOffer, offerFarID, "Far", -64.41, -31.21, userID).Error; err != nil {
		t.Fatalf("insert far offer: %v", err)
	}

	offerRepo := NewOfferRepository(db)
	offers, err := offerRepo.FindNearby(-31.2, -64.4, 5000)
	if err != nil {
		t.Fatalf("FindNearby() error = %v", err)
	}
	if len(offers) != 2 || offers[0].ID != offerNearID || offers[0].PostedBy.ID != userID.String() {
		t.Fatalf("FindNearby() = %#v", offers)
	}
	if _, err := offerRepo.FindByID(uuid.New()); err != gorm.ErrRecordNotFound {
		t.Fatalf("FindByID(missing) error = %v, want record not found", err)
	}

	commentRepo := NewCommentRepository(db)
	comment := &models.Comment{ID: uuid.New(), OfferID: offerNearID, UserID: userID, Body: "Available"}
	if err := commentRepo.Create(comment); err != nil {
		t.Fatalf("create comment: %v", err)
	}
	comments, err := commentRepo.FindByOfferID(offerNearID)
	if err != nil || len(comments) != 1 || comments[0].PostedBy.ID != userID.String() {
		t.Fatalf("FindByOfferID() = %#v, %v", comments, err)
	}
	found, err := offerRepo.FindByID(offerNearID)
	if err != nil || found.CommentsCount != 1 {
		t.Fatalf("offer comments count = %d, %v", found.CommentsCount, err)
	}

	voteRepo := NewOfferVoteRepository(db)
	confirmations, invalidations, err := voteRepo.Vote(offerNearID, userID, models.VoteTypeValidate)
	if err != nil || confirmations != 1 || invalidations != 0 {
		t.Fatalf("validate counts = %d/%d, %v", confirmations, invalidations, err)
	}
	confirmations, invalidations, err = voteRepo.Vote(offerNearID, userID, models.VoteTypeInvalidate)
	if err != nil || confirmations != 0 || invalidations != 1 {
		t.Fatalf("switched counts = %d/%d, %v", confirmations, invalidations, err)
	}

	if err := db.Exec("ALTER TABLE offers ADD CONSTRAINT force_vote_rollback CHECK (confirmations_count = 0)").Error; err != nil {
		t.Fatalf("add rollback constraint: %v", err)
	}
	if _, _, err := voteRepo.Vote(offerNearID, userID, models.VoteTypeValidate); err == nil {
		t.Fatal("Vote() error = nil, want forced transaction failure")
	}
	var storedType models.VoteType
	if err := db.Model(&models.OfferVote{}).Select("type").Where("offer_id = ? AND user_id = ?", offerNearID, userID).Scan(&storedType).Error; err != nil {
		t.Fatalf("read vote after rollback: %v", err)
	}
	if storedType != models.VoteTypeInvalidate {
		t.Fatalf("vote after rollback = %q, want invalidate", storedType)
	}
}
