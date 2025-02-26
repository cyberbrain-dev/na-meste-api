package authentication

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Contains payload of a JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Generates a JWT with the role encoded
func GenerateJWT(id uint, role string) (string, error) {
	// creating a payload
	claims := Claims{
		UserID: id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "na-meste-api",
		},
	}

	// creating an unsigned token
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return unsignedToken.SignedString("my-secret-key")
}
