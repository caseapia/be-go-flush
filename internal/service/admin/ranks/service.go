package adminRanks

import (
	AdminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	"github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/caseapia/goproject-flush/internal/service/logger"
)

type RanksService struct {
	repo           *AdminRanksRepository.RanksRepository
	userRankSetter contracts.UserRankSetter
	logger         *logger.LoggerService
}

func NewRanksService(r *AdminRanksRepository.RanksRepository, u contracts.UserRankSetter) *RanksService {
	return &RanksService{
		repo:           r,
		userRankSetter: u,
	}
}

func (s *RanksService) SetUserRankSetter(u contracts.UserRankSetter) {
	s.userRankSetter = u
}
