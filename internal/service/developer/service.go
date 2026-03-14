package developer

import (
	"context"
	"fmt"

	"github.com/caseapia/goproject-flush/config"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	"github.com/gookit/slog"
)

type Service struct {
	repo    mysql.Repository
	logger  logger.Service
	manager *config.ServiceManager
}

func NewService(repo mysql.Repository, logger logger.Service, manager *config.ServiceManager) *Service {
	return &Service{repo: repo, logger: logger, manager: manager}
}

func (s *Service) PopulateServices(ctx context.Context) map[string]*config.ServiceStatus {
	statuses := s.manager.PopulateServices()

	return statuses
}

func (s *Service) ServiceInteraction(ctx context.Context, name string, updatedBy uint64, action enums.ServiceAction) (bool, error) {
	user, err := s.repo.SearchByID(ctx, updatedBy)
	if err != nil {
		return false, err
	}

	switch action {
	case enums.Enable:
		enabledState := true
		err = s.manager.SetServiceEnabled(name, enabledState, user.Name)
		if err != nil {
			return false, err
		}

		s.logger.Log(ctx, enums.StaffCommonLogger, &updatedBy, nil, enums.EnableService, fmt.Sprintf("Service: %s", name))

		return true, err
	case enums.Disable:
		enabledState := false
		err = s.manager.SetServiceEnabled(name, enabledState, user.Name)
		if err != nil {
			return false, err
		}

		s.logger.Log(ctx, enums.StaffCommonLogger, &updatedBy, nil, enums.DisableService, fmt.Sprintf("Service: %s", name))

		return false, err
	case enums.Status:
		state := s.manager.IsServiceEnabled(name)

		return state, nil
	default:
		slog.Error("Unknown action. Use: enable|disable|status")
	}

	return false, err
}
