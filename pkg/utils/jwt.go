package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("your_secret_key") // TODO: Use environment variable

type Claims struct {
	Username string    `json:"username"`
	Role     string    `json:"role"`
	OutletID *uint     `json:"outlet_id,omitempty"`
	ID       uuid.UUID `json:"id"` // User's ExternalID
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token.
func GenerateToken(username, role string, outletID *uint, userUuid uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role,
		OutletID: outletID,
		ID:       userUuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken parses and validates a JWT token.
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
