package adminRanks

import (
	AdminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	"github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/caseapia/goproject-flush/internal/service/logger"
)

type RanksService struct {
	ranksRepo      *AdminRanksRepository.RanksRepository
	userRankSetter contracts.UserRankSetter
	logger         *logger.LoggerService
}

func NewRanksService(
	ranksRepo *AdminRanksRepository.RanksRepository,
	userRankSetter contracts.UserRankSetter,
	logger *logger.LoggerService,
) *RanksService {
	return &RanksService{
		ranksRepo:      ranksRepo,
		userRankSetter: userRankSetter,
		logger:         logger,
	}
}

func (s *RanksService) SetUserRankSetter(u contracts.UserRankSetter) {
	s.userRankSetter = u
}
