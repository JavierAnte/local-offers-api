package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const tokenTTL = 30 * 24 * time.Hour

type claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, userID uuid.UUID) (string, error) {
	c := claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret string, tokenString string) (uuid.UUID, error) {
	c := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid or expired token")
	}

	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return uuid.Nil, errors.New("invalid token subject")
	}

	return userID, nil
}
