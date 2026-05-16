package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("CHANGE_THIS_TO_A_LONG_RANDOM_SECRET")

type Claims struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	IsPremium bool   `json:"is_premium"`
	jwt.RegisteredClaims
}

func GenerateJWT(
	userID int,
	username string,
	isPremium bool,
) (string, error) {

	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserID:    userID,
		Username:  username,
		IsPremium: isPremium,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(JwtSecret)
}