package rankshandler

import ranksservice "github.com/caseapia/goproject-flush/internal/service/admin/ranks"

type RanksHandler struct {
	service *ranksservice.RanksService
}

func NewRanksService(s *ranksservice.RanksService) *RanksHandler {
	return &RanksHandler{
		service: s,
	}
}
