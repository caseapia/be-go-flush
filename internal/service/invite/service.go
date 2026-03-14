package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/caseapia/goproject-flush/config"
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	inviteutils "github.com/caseapia/goproject-flush/pkg/utils/invite"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Service struct {
	repo           mysql.Repository
	logger         logger.Service
	serviceManager *config.ServiceManager
}

func NewService(inviteRepo mysql.Repository, logger logger.Service, serviceManager *config.ServiceManager) *Service {
	return &Service{repo: inviteRepo, logger: logger, serviceManager: serviceManager}
}

func (s *Service) GetInviteCodes(ctx context.Context) ([]models.Invite, error) {
	if s.serviceManager.IsServiceEnabled("invites") != true {
		return nil, fiber.NewError(403, "service disabled by an admin")
	}

	invites, err := s.repo.SearchAllInvites(ctx)
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
		}).Error("error when fetching invite codes")
		return nil, err
	}

	return invites, nil
}

func (s *Service) GetInviteByID(ctx context.Context, inviteID string) (*models.Invite, error) {
	if s.serviceManager.IsServiceEnabled("invites") != true {
		return nil, fiber.NewError(403, "service disabled by an admin")
	}

	inviteInfo, err := s.repo.SearchInviteByCode(ctx, inviteID)
	if err != nil {
		return nil, fiber.NewError(500, err.Error())
	}

	return inviteInfo, nil
}

func (s *Service) CreateInvite(ctx context.Context, createdBy uint64) (*models.Invite, error) {
	if s.serviceManager.IsServiceEnabled("invites") != true {
		return nil, fiber.NewError(403, "service disabled by an admin")
	}

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

	if err := s.repo.CreateInvite(ctx, s.repo.DB, invite); err != nil {
		return nil, err
	}

	addInfo := fmt.Sprintf("ID: %v\nCode: %s", invite.ID, invite.Code)
	s.logger.Log(ctx, enums.StaffCommonLogger, &createdBy, nil, enums.CreateInvite, addInfo)

	return invite, nil
}

func (s *Service) ValidateInvite(ctx context.Context, code string) (*models.Invite, error) {
	if s.serviceManager.IsServiceEnabled("invites") != true {
		return nil, fiber.NewError(403, "service disabled by an admin")
	}

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
	if s.serviceManager.IsServiceEnabled("invites") != true {
		return fiber.NewError(403, "service disabled by an admin")
	}

	invite, err := s.repo.SearchInviteByCode(ctx, code)
	if err != nil {
		return err
	}

	if invite.Used {
		return &fiber.Error{Code: 403, Message: "invite already used"}
	}

	return s.repo.MarkInviteAsUsed(ctx, s.repo.DB, invite.ID, userID)
}

func (s *Service) DeleteInvite(ctx context.Context, adminID uint64, inviteID uint64) error {
	if s.serviceManager.IsServiceEnabled("invites") != true {
		return fiber.NewError(403, "service disabled by an admin")
	}

	oldInvite, err := s.repo.SearchInviteByID(ctx, inviteID)
	if err != nil {
		return err
	}

	newErr := s.repo.DeleteInvite(ctx, s.repo.DB, inviteID)
	if newErr != nil {
		return newErr
	}

	addInfo := fmt.Sprintf("ID: %v\nCode: %s", inviteID, oldInvite.Code)
	s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, nil, enums.DeleteInvite, addInfo)

	return newErr
}
