package AdminModule

import (
	AdminRanksHandler "github.com/caseapia/goproject-flush/internal/handler/admin/ranks"
	AdminUserHandler "github.com/caseapia/goproject-flush/internal/handler/admin/user"
	AdminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	AdminRanksService "github.com/caseapia/goproject-flush/internal/service/admin/ranks"
	Contracts "github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type AdminModule struct {
	RanksHandler *AdminRanksHandler.RanksHandler
	RanksService *AdminRanksService.RanksService
	UserHandler  *AdminUserHandler.AdminUserHandler
}

func NewAdminModule(db *bun.DB, userRankSetter Contracts.UserRankSetter, userHandler *AdminUserHandler.AdminUserHandler) *AdminModule {
	ranksRepo := AdminRanksRepository.NewRanksRepository(db)
	ranksSrv := AdminRanksService.NewRanksService(ranksRepo, userRankSetter)
	ranksHandler := AdminRanksHandler.NewRanksHandler(ranksSrv)

	return &AdminModule{
		RanksHandler: ranksHandler,
		RanksService: ranksSrv,
		UserHandler:  userHandler,
	}
}

func (m *AdminModule) RegisterRoutes(app fiber.Router) {
	admin := app.Group("/admin")

	// Ranks
	admin.Get("/ranks", m.RanksHandler.GetRanksList)
	admin.Post("/setstaff/:id", m.RanksHandler.SetStaffRank)
	admin.Post("/setdeveloper/:id", m.RanksHandler.SetDeveloperRank)

	// User actions
	if m.UserHandler != nil {
		admin.Delete("/delete/:id", m.UserHandler.DeleteUser)
		admin.Put("/restore/:id", m.UserHandler.RestoreUser)
		admin.Put("/create", m.UserHandler.CreateUser)
		admin.Patch("/ban/:id", m.UserHandler.BanUser)
		admin.Delete("/unban/:id", m.UserHandler.UnbanUser)
	}
}
