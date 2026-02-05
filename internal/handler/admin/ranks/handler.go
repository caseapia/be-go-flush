package AdminRanksHandler

import ranksservice "github.com/caseapia/goproject-flush/internal/service/admin/ranks"

type RanksHandler struct {
	service *ranksservice.RanksService
}

func NewRanksHandler(s *ranksservice.RanksService) *RanksHandler {
	return &RanksHandler{
		service: s,
	}
}
