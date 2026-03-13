package logger

import (
	"context"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/pkg/utils/enums"
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
	loggerType enums.LoggerType,
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

	act, ok := action.(enums.LoggerAction)
	if !ok {
		slog.Error("invalid action type")
		return
	}
	base.Action = string(act)

	switch loggerType {
	case enums.StaffPunishmentLogger:
		err := s.repo.SaveLog(ctx, s.repo.DB, &models.PunishmentLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})
		if err != nil {
			slog.WithData(slog.M{
				"error": err,
				"base":  base,
				"type:": loggerType,
			}).Error("error in log saving")
		}

	case enums.StaffCommonLogger:
		err := s.repo.SaveLog(ctx, s.repo.DB, &models.CommonLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})
		if err != nil {
			slog.WithData(slog.M{
				"error": err,
				"base":  base,
				"type:": loggerType,
			}).Error("error in log saving")
		}

	case enums.AdminTicketLogger:
		err := s.repo.SaveLog(ctx, s.repo.DB, &models.TicketsLog{
			BaseLog: base,
			AdminID: *adminID,
			UserID:  userID,
		})
		if err != nil {
			slog.WithData(slog.M{
				"error": err,
				"base":  base,
				"type:": loggerType,
			}).Error("error in log saving")
		}

	case enums.AdminAuthLogger:
		err := s.repo.SaveLog(ctx, s.repo.DB, &models.AuthLog{
			BaseLog: base,
			UserID:  userID,
		})
		if err != nil {
			slog.WithData(slog.M{
				"error": err,
				"base":  base,
				"type:": loggerType,
			}).Error("error in log saving")
		}

	default:
		slog.WithData(slog.M{
			"base": base,
			"type": loggerType,
		}).Error("unknown logger type")
	}
}
