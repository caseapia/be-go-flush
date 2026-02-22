package logger

import (
	"context"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/gookit/slog"
)

type Service struct {
	repo mysql.Repository
}

func NewService(r mysql.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetCommonStaffLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.CommonLog, int, error) {
	return s.repo.GetCommonLogs(ctx, startDate, endDate, keywords)
}

func (s *Service) GetPunishmentStaffLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.PunishmentLog, int, error) {
	return s.repo.GetPunishmentLogs(ctx, startDate, endDate, keywords)
}

func (s *Service) GetTicketsAdminLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.TicketsLog, int, error) {
	return s.repo.GetTicketsLog(ctx, startDate, endDate, keywords)
}

func (s *Service) GetAuthAdminLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.AuthLog, int, error) {
	return s.repo.GetAuthLogs(ctx, startDate, endDate, keywords)
}

func (s *Service) Log(
	ctx context.Context,
	loggerType models.LoggerType,
	adminID *uint64,
	userID *uint64,
	action interface{},
	additional ...string,
) {
	var addInfo *string
	if len(additional) > 0 {
		addInfo = &additional[0]
	}

	base := models.BaseLog{
		AdditionalInfo: addInfo,
		Date:           time.Now(),
	}

	act, ok := action.(models.Action)
	if !ok {
		slog.Error("invalid action type")
		return
	}
	base.Action = act

	switch loggerType {

	case models.StaffPunishmentLogger:
		s.repo.SaveLog(ctx, &models.PunishmentLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})

	case models.StaffCommonLogger:
		s.repo.SaveLog(ctx, &models.CommonLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})

	case models.AdminTicketLogger:
		s.repo.SaveLog(ctx, &models.TicketsLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})

	case models.AdminAuthLogger:
		s.repo.SaveLog(ctx, &models.AuthLog{
			BaseLog: base,
			UserID:  userID,
		})

	default:
		slog.Error("unknown logger type")
	}
}
