package adminRanks

import adminRanks "github.com/caseapia/goproject-flush/internal/service/admin/ranks"

type Handler struct {
	service *adminRanks.RanksService
}

func NewHandler(s *adminRanks.RanksService) *Handler {
	return &Handler{
		service: s,
	}
}
