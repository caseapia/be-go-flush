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
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseAccessToken(tokenStr string) (*Claims, error) {
	if len(tokenStr) == 0 {
		return nil, jwt.ErrTokenMalformed
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func GetJWTSecret() []byte {
	return jwtSecret
}
