package developer

import (
	"bytes"
	"context"
	"fmt"
	"runtime/pprof"
	"time"

	"github.com/caseapia/goproject-flush/config"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
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

func (s *Service) DebugStack(c *fiber.Ctx) (string, error) {
	var buf bytes.Buffer

	pprof.Lookup("goroutine").WriteTo(&buf, 2)

	return buf.String(), nil
}

func (s *Service) ServerInfo(c *fiber.Ctx) (*int64, *float64, *float64, *uint64, *uint64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	cpuPercent, err := cpu.Percent(time.Millisecond*100, false)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	uptime, err := host.Uptime()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var cpuUsage float64
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	now := time.Now().Unix()
	usedGB := v.Used / 1024 / 1024 / 1024

	return &now, &cpuUsage, &v.UsedPercent, &usedGB, &uptime, nil
}

func (s *Service) Ping(c *fiber.Ctx) string {
	return "ok"
}
