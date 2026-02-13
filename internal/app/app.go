package app

import (
	"github.com/caseapia/goproject-flush/config"
	"github.com/caseapia/goproject-flush/internal/db"
	adminUserHandler "github.com/caseapia/goproject-flush/internal/handler/admin/user"
	authHandler "github.com/caseapia/goproject-flush/internal/handler/user/auth"
	"github.com/caseapia/goproject-flush/internal/middleware"
	adminInvite "github.com/caseapia/goproject-flush/internal/service/admin/invite"
	adminRanksService "github.com/caseapia/goproject-flush/internal/service/admin/ranks"
	adminUserService "github.com/caseapia/goproject-flush/internal/service/admin/user"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	mysqlRepo "github.com/caseapia/goproject-flush/internal/repository/mysql"
	adminModule "github.com/caseapia/goproject-flush/pkg/module/admin"
	loggerModule "github.com/caseapia/goproject-flush/pkg/module/logger"
	userModule "github.com/caseapia/goproject-flush/pkg/module/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"

)

func NewApp() (*fiber.App, error) {
	// Загружаем конфиг
	config.LoadEnv()

	// Подключаем БД
	dbs, err := db.NewDatabases()
	if err != nil {
		return nil, err
	}

	// Настройка логгера
	slog.Configure(func(l *slog.SugaredLogger) {
		if f, ok := l.Formatter.(*slog.TextFormatter); ok {
			f.EnableColor = true
		}
	})
	slog.SetFormatter(slog.NewJSONFormatter())

	// Универсальный репозиторий
	repo := mysqlRepo.NewRepository(dbs.Main)

	// Сервисы
	loggerSrv := loggerService.NewLoggerService(repo)
	userRankSetter := adminRanksService.NewUserRankSetter(repo, loggerSrv)
	ranksSrv := adminRanksService.NewRanksService(repo, userRankSetter, loggerSrv)
	adminUserSrv := adminUserService.NewAdminUserService(repo, ranksSrv, loggerSrv)
	inviteSrv := adminInvite.NewAdminService(repo, ranksSrv)
	authSrv := authHandler.NewService(repo, inviteSrv)

	// Хэндлеры
	authH := authHandler.NewHandler(authSrv)
	adminUserH := adminUserHandler.NewAdminUserHandler(adminUserSrv)

	// Модули
	userM := userModule.NewUserModule(repo, loggerSrv, ranksSrv)
	adminM := adminModule.NewAdminModule(repo, ranksSrv, adminUserH, loggerSrv, inviteSrv)
	loggerM := loggerModule.NewLoggerModule(repo, loggerSrv)

	// Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.AppErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://fe-go-flush.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Роуты
	api := app.Group("/api")
	userM.RegisterRoutes(api)
	adminM.RegisterRoutes(api)
	loggerM.RegisterRoutes(api)
	authH.RegisterRoutes(api)

	return app, nil
}
