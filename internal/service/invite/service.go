package invite

import (
	"context"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	inviteutils "github.com/caseapia/goproject-flush/pkg/utils/invite"
	"github.com/gofiber/fiber/v2"
)

type InviteRepository interface {
	CreateInvite(ctx context.Context, invite *models.Invite) error
	DeleteInvite(ctx context.Context, inviteID uint64) error
	SearchInviteByCode(ctx context.Context, code string) (*models.Invite, error)
	MarkInviteAsUsed(ctx context.Context, inviteID, usedBy uint64) error
}

type Service struct {
	repo InviteRepository
}

func NewService(repo InviteRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateInvite(ctx context.Context, createdBy uint64) (*models.Invite, error) {
	code, err := inviteutils.GenerateCode()
	if err != nil {
		return nil, err
	}

	invite := &models.Invite{
		Code:      code,
		CreatedBy: createdBy,
		Used:      false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateInvite(ctx, invite); err != nil {
		return nil, err
	}

	return invite, nil
}

func (s *Service) ValidateInvite(ctx context.Context, code string) (*models.Invite, error) {
	invite, err := s.repo.SearchInviteByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if invite.Used {
		return nil, &fiber.Error{Code: 403, Message: "invite already used"}
	}

	return invite, nil
}

func (s *Service) UseInvite(ctx context.Context, code string, userID uint64) error {
	invite, err := s.repo.SearchInviteByCode(ctx, code)
	if err != nil {
		return err
	}

	if invite.Used {
		return &fiber.Error{Code: 403, Message: "invite already used"}
	}

	return s.repo.MarkInviteAsUsed(ctx, invite.ID, userID)
}

func (s *Service) DeleteInvite(ctx context.Context, inviteID uint64) error {
	return s.repo.DeleteInvite(ctx, inviteID)
}
