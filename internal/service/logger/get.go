package logger

import (
	"context"

	models "github.com/caseapia/goproject-flush/internal/models/logger"
)

func (s *LoggerService) GetAllLogs(ctx context.Context) ([]models.BaseLog, error) {
	return s.repo.GetAllLogs(ctx)
}

func (s *LoggerService) GetCommonLogs(ctx context.Context) ([]models.CommonLog, error) {
	return s.repo.GetCommonLogs(ctx)
}

func (s *LoggerService) GetPunishmentLogs(ctx context.Context) ([]models.PunishmentLog, error) {
	return s.repo.GetPunishmentLogs(ctx)
}
