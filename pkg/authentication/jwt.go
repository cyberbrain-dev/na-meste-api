package authentication

import (
	"fmt"
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

	return unsignedToken.SignedString([]byte("my-secret-key"))
}

// Parses and verifies the JWT
func ParseJWT(tokenString string) (*Claims, error) {
	// parsing the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte("my-secret-key"), nil
		},
	)
	// if smth went wrong
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// getting the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid or expired token")
}
