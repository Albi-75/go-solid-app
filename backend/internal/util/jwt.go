package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your-secret-key") // use env variable in prod

func GenerateJWT(userID int32, username string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        string(userID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
