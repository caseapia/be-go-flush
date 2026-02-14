package invite

import (
	"context"
	"strconv"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	inviteutils "github.com/caseapia/goproject-flush/pkg/utils/invite"
	"github.com/gofiber/fiber/v2"
)

type InviteRepository interface {
	SearchAllInvites(ctx context.Context) ([]models.InviteDTO, error)
	CreateInvite(ctx context.Context, invite *models.Invite) error
	DeleteInvite(ctx context.Context, inviteID uint64) error
	SearchInviteByCode(ctx context.Context, code string) (*models.Invite, error)
	MarkInviteAsUsed(ctx context.Context, inviteID, usedBy uint64) error
}

type Service struct {
	inviteRepo InviteRepository
	logger     logger.Service
}

func NewService(inviteRepo InviteRepository, logger logger.Service) *Service {
	return &Service{inviteRepo: inviteRepo, logger: logger}
}

func (s *Service) GetInviteCodes(ctx context.Context) ([]models.InviteDTO, error) {
	return s.inviteRepo.SearchAllInvites(ctx)
}

func (s *Service) GetInviteByID(ctx context.Context, inviteID string) (*models.Invite, error) {
	inviteInfo, err := s.inviteRepo.SearchInviteByCode(ctx, inviteID)
	if err != nil {
		return nil, fiber.NewError(500, err.Error())
	}

	return inviteInfo, nil
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

	if err := s.inviteRepo.CreateInvite(ctx, invite); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx, models.CommonLogger, createdBy, nil, models.CreateInvite, invite.Code)

	return invite, nil
}

func (s *Service) ValidateInvite(ctx context.Context, code string) (*models.Invite, error) {
	invite, err := s.inviteRepo.SearchInviteByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if invite.Used {
		return nil, &fiber.Error{Code: 403, Message: "invite already used"}
	}

	return invite, nil
}

func (s *Service) UseInvite(ctx context.Context, code string, userID uint64) error {
	invite, err := s.inviteRepo.SearchInviteByCode(ctx, code)
	if err != nil {
		return err
	}

	if invite.Used {
		return &fiber.Error{Code: 403, Message: "invite already used"}
	}

	return s.inviteRepo.MarkInviteAsUsed(ctx, invite.ID, userID)
}

func (s *Service) DeleteInvite(ctx context.Context, adminID uint64, inviteID uint64) error {
	_ = s.logger.Log(ctx, models.CommonLogger, adminID, nil, models.DeleteInvite, strconv.FormatUint(inviteID, 10))

	return s.inviteRepo.DeleteInvite(ctx, inviteID)
}
