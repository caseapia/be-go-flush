package admin

import (
	rankshandler "github.com/caseapia/goproject-flush/internal/handler/admin/ranks"
	ranksrepo "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	ranksservice "github.com/caseapia/goproject-flush/internal/service/admin/ranks"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type AdminModule struct {
	RanksHandler *rankshandler.RanksHandler
}

func NewAdminModule(db *bun.DB) *AdminModule {
	ranksRepo := ranksrepo.NewRanksRepository(db)

	ranksSrv := ranksservice.NewRankService(ranksRepo)

	ranksHandler := rankshandler.NewRanksService(ranksSrv)

	return &AdminModule{
		RanksHandler: ranksHandler,
	}
}

func (m *AdminModule) RegisterRoutes(app fiber.Router) {
	admin := app.Group("/admin")

	// Ranks
	admin.Get("/ranks", m.RanksHandler.GetRanksList)
}
