package app

import (
	"github.com/caseapia/goproject-flush/config"
	database "github.com/caseapia/goproject-flush/internal/db"
	"github.com/caseapia/goproject-flush/internal/handler/auth"
	"github.com/caseapia/goproject-flush/internal/handler/invite"
	"github.com/caseapia/goproject-flush/internal/handler/logger"
	"github.com/caseapia/goproject-flush/internal/handler/ranks"
	"github.com/caseapia/goproject-flush/internal/handler/user"
	mysqlRepo "github.com/caseapia/goproject-flush/internal/repository/mysql"
	authService "github.com/caseapia/goproject-flush/internal/service/auth"
	inviteService "github.com/caseapia/goproject-flush/internal/service/invite"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	ranksService "github.com/caseapia/goproject-flush/internal/service/ranks"
	userService "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
)

func NewApp() (*fiber.App, error) {
	config.LoadEnv()
	setupLogger()

	dbs, err := database.NewDatabases()
	if err != nil {
		return nil, err
	}

	mainRepo := mysqlRepo.NewRepository(dbs.Main)
	logsRepo := mysqlRepo.NewRepository(dbs.Logs)

	loggerSrv := loggerService.NewService(*logsRepo)

	ranksSrv := ranksService.NewService(mainRepo, loggerSrv)

	userSrv := userService.NewService(mainRepo, loggerSrv)

	inviteSrv := inviteService.NewService(mainRepo)

	authSrv := authService.NewService(*mainRepo)

	handlers := struct {
		auth   *auth.Handler
		user   *user.Handler
		invite *invite.Handler
		logger *logger.Handler
		ranks  *ranks.Handler
	}{
		auth:   auth.NewHandler(authSrv, inviteSrv),
		user:   user.NewUserHandler(userSrv, ranksSrv),
		invite: invite.NewHandler(inviteSrv),
		logger: logger.NewHandler(loggerSrv),
		ranks:  ranks.NewHandler(ranksSrv),
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://fe-go-flush.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api := app.Group("/api")

	// ! Public routes
	handlers.auth.RegisterRoutes(api)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// ! Private Routes
	private := api.Group("")
	private.Use(auth.AuthMiddleware(authSrv))

	handlers.user.RegisterRoutes(private)
	handlers.invite.RegisterRoutes(private)
	handlers.logger.RegisterRoutes(private)
	handlers.ranks.RegisterRoutes(private)
	handlers.auth.RegisterPrivateRoute(private)

	return app, nil
}

func setupLogger() {
	slog.Configure(func(l *slog.SugaredLogger) {
		if f, ok := l.Formatter.(*slog.TextFormatter); ok {
			f.EnableColor = true
		}
	})
	slog.SetFormatter(slog.NewJSONFormatter())
}
