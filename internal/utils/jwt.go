package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key")

type Claims struct {
	UserID    uint64 `json:"sub"`
	SessionID string `json:"sid"`
	TokenVer  int    `json:"tv"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint64, sessionID string, tokenVer int) (string, error) {
	claims := Claims{
		UserID:    userID,
		SessionID: sessionID,
		TokenVer:  tokenVer,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
