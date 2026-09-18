package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateAndParseToken(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	token, err := GenerateToken("test-secret", userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	got, err := ParseToken("test-secret", token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if got != userID {
		t.Fatalf("ParseToken() user = %v, want %v", got, userID)
	}
}

func TestParseTokenRejectsInvalidTokens(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	validToken, err := GenerateToken("correct-secret", userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	expiredClaims := claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		},
	}
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).
		SignedString([]byte("correct-secret"))
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}

	invalidSubjectClaims := claims{
		UserID: "not-a-uuid",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	invalidSubjectToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidSubjectClaims).
		SignedString([]byte("correct-secret"))
	if err != nil {
		t.Fatalf("sign invalid-subject token: %v", err)
	}
	wrongAlgorithmToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, claims{
		UserID:           userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	}).SignedString([]byte("correct-secret"))
	if err != nil {
		t.Fatalf("sign wrong-algorithm token: %v", err)
	}

	tests := []struct {
		name   string
		secret string
		token  string
	}{
		{name: "malformed", secret: "correct-secret", token: "not-a-token"},
		{name: "wrong secret", secret: "wrong-secret", token: validToken},
		{name: "expired", secret: "correct-secret", token: expiredToken},
		{name: "invalid subject", secret: "correct-secret", token: invalidSubjectToken},
		{name: "wrong algorithm", secret: "correct-secret", token: wrongAlgorithmToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := ParseToken(tt.secret, tt.token); err == nil {
				t.Fatal("ParseToken() error = nil, want error")
			}
		})
	}
}
