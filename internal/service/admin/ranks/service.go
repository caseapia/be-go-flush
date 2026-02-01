package ranksservice

import ranksrepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"

type RanksService struct {
	repo *ranksrepository.RanksRepository
}

func NewRankService(r *ranksrepository.RanksRepository) *RanksService {
	return &RanksService{
		repo: r,
	}
}
