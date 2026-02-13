package auth

import (
	"context"
	"errors"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/hash"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

)

type Service struct {
	repository mysql.Repository
}

func NewService(userRepo mysql.Repository) *Service {
	return &Service{repository: userRepo}
}

func (s *Service) Register(ctx context.Context, name, invite, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		Password:     string(hash),
		TokenVersion: 1,
		IsVerified:   true,
	}

	return s.repository.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, login, password, userAgent, ip string) (string, string, error) {
	user, err := s.repository.GetByEmail(ctx, login)
	if err != nil {
		return "", "", &fiber.Error{Code: 404, Message: "invalid credentials"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", &fiber.Error{Code: 404, Message: "invalid credentials"}
	}

	sessionID := uuid.NewString()
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	refreshHash := hash.HashToken(refreshToken)

	session := &models.Session{
		ID:          sessionID,
		UserID:      user.ID,
		RefreshHash: refreshHash,
		UserAgent:   userAgent,
		IPLast:      ip,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:   time.Now(),
	}

	if err := s.repository.CreateSession(ctx, session); err != nil {
		return "", "", err
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, sessionID, user.TokenVersion)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	refreshHash := hash.HashToken(refreshToken)

	session, err := s.repository.GetSessionByHash(ctx, refreshHash)
	if err != nil || session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("invalid or expired session")
	}

	user, err := s.repository.SearchUserByID(ctx, session.UserID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	session.RefreshHash = hash.HashToken(newRefreshToken)
	session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	if err := s.repository.UpdateSession(ctx, session); err != nil {
		return "", "", err
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, session.ID, user.TokenVersion)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	return s.repository.RevokeSession(ctx, sessionID)
}
