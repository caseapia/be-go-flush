package auth

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key")

type Claims struct {
	UserID    uint64 `json:"sub"`
	SessionID string `json:"sid"`
	TokenVer  int    `json:"tv"`
	jwt.RegisteredClaims
}

func (s *Service) ValidateAccessToken(tokenStr string) (uint64, error) {
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &fiber.Error{Code: 400, Message: "unexpected signing method"}
		}
		return jwtSecret, nil
	})
	if err != nil {
		return 0, &fiber.Error{Code: 404, Message: "invalid token"}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, &fiber.Error{Code: 404, Message: "invalid token claims"}
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return 0, &fiber.Error{Code: 404, Message: "token expired"}
	}

	return claims.UserID, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
