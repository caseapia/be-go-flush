package usermodule

import (
	userhandler "github.com/caseapia/goproject-flush/internal/handler/user"
	userrepo "github.com/caseapia/goproject-flush/internal/repository/user"
	userservice "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type UserModule struct {
	Handler *userhandler.UserHandler
}

func NewUserModule(db *bun.DB) *UserModule {
	repo := userrepo.NewUserRepository(db)
	srv := userservice.NewUserService(repo, nil)
	h := userhandler.NewUserHandler(srv)

	return &UserModule{Handler: h}
}

func (m *UserModule) RegisterRoutes(app fiber.Router) {
	userhandler.RegisterRoutes(app, m.Handler)
}
