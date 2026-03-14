package app

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	"github.com/caseapia/goproject-flush/config"
	"github.com/caseapia/goproject-flush/internal/clients/discord"
	database "github.com/caseapia/goproject-flush/internal/db"
	"github.com/caseapia/goproject-flush/internal/handler/auth"
	"github.com/caseapia/goproject-flush/internal/handler/badges"
	"github.com/caseapia/goproject-flush/internal/handler/changelog"
	"github.com/caseapia/goproject-flush/internal/handler/developer"
	"github.com/caseapia/goproject-flush/internal/handler/invite"
	"github.com/caseapia/goproject-flush/internal/handler/logger"
	"github.com/caseapia/goproject-flush/internal/handler/notifications"
	"github.com/caseapia/goproject-flush/internal/handler/ranks"
	"github.com/caseapia/goproject-flush/internal/handler/tickets"
	"github.com/caseapia/goproject-flush/internal/handler/user"
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	mysqlRepo "github.com/caseapia/goproject-flush/internal/repository/mysql"
	authService "github.com/caseapia/goproject-flush/internal/service/auth"
	badgesService "github.com/caseapia/goproject-flush/internal/service/badges"
	changelogService "github.com/caseapia/goproject-flush/internal/service/changelog"
	developerService "github.com/caseapia/goproject-flush/internal/service/developer"
	inviteService "github.com/caseapia/goproject-flush/internal/service/invite"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	notifyService "github.com/caseapia/goproject-flush/internal/service/notifications"
	ranksService "github.com/caseapia/goproject-flush/internal/service/ranks"
	ticketsService "github.com/caseapia/goproject-flush/internal/service/tickets"
	userService "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func NewApp() (*fiber.App, error) {
	migrateFlag := flag.Bool("migrate", false, "run database migrations and exit")
	debugFlag := flag.Bool("debug", false, "debug app with display of incoming requests")
	flag.Parse()

	config.LoadEnv()
	setupLogger()

	cfg := config.Load()

	discordClient := discord.NewClient(
		cfg.DiscordClientID,
		cfg.DiscordClientSecret,
		cfg.DiscordRedirectURI(),
	)

	dbs, err := database.NewDatabases()
	if err != nil {
		return nil, err
	}
	mainRepo := mysqlRepo.NewRepository(dbs.Main)
	logsRepo := mysqlRepo.NewRepository(dbs.Logs)

	if *migrateFlag {
		if err := mysql.RunMigrations(dbs.Main, mysql.MainModels); err != nil {
			log.Fatal(err)
		}
		if err := mysql.RunMigrations(dbs.Logs, mysql.LogModels); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	manager, err := config.NewServiceManager("services.json")
	if err != nil {
		return nil, err
	}

	loggerSrv := loggerService.NewService(*logsRepo)
	badgesSrv := badgesService.NewService(*mainRepo, *loggerSrv)
	notifySrv := notifyService.NewService(*mainRepo, *loggerSrv, manager)
	ranksSrv := ranksService.NewService(*mainRepo, *loggerSrv)
	userSrv := userService.NewService(*mainRepo, *loggerSrv, *notifySrv)
	inviteSrv := inviteService.NewService(*mainRepo, *loggerSrv, manager)
	authSrv := authService.NewService(*mainRepo, *loggerSrv, inviteSrv, *notifySrv, discordClient, cfg, manager)
	ticketsSrv := ticketsService.NewService(*mainRepo, *notifySrv, *loggerSrv, manager)
	changelogSrv := changelogService.NewService(*mainRepo, *loggerSrv, *notifySrv, manager)
	developerSrv := developerService.NewService(*mainRepo, *loggerSrv, manager)

	badgesHandler := badges.NewHandler(badgesSrv)
	authHandler := auth.NewHandler(authSrv)
	userHandler := user.NewUserHandler(userSrv, ranksSrv, authSrv)
	inviteHandler := invite.NewHandler(inviteSrv)
	loggerHandler := logger.NewHandler(loggerSrv)
	ranksHandler := ranks.NewHandler(ranksSrv)
	notifyHandler := notifications.NewHandler(notifySrv)
	ticketsHandler := tickets.NewHandler(ticketsSrv)
	changelogHandler := changelog.NewHandler(changelogSrv)
	developerHandler := developer.NewHandler(developerSrv)

	events := manager.Subscribe()
	go func() {
		for event := range events {
			slog.Infof("Service %s is now %v", event.ServiceName, event.Enabled)
		}
	}()

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		Concurrency:  256 * 1024,
		ProxyHeader:  fiber.HeaderXForwardedFor,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
					"code":  e.Code,
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": err.Error(),
			})
		},
	})

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	app.Get("/debug/stack", func(c *fiber.Ctx) error {
		var buf bytes.Buffer
		pprof.Lookup("goroutine").WriteTo(&buf, 2)
		return c.Type("text/plain").Send(buf.Bytes())
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8080, https://fe-go-flush.vercel.app, https://dash.dontkillme.lol",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cache-Control",
		AllowCredentials: true,
	}))

	if *debugFlag {
		fmt.Print("debug started. Incoming requests will be displayed")

		app.Use(func(c *fiber.Ctx) error {
			start := time.Now()

			err := c.Next()

			stop := time.Since(start)

			slog.WithData(slog.M{
				"method":  c.Method(),
				"path":    c.Path(),
				"status":  c.Response().StatusCode(),
				"latency": stop.String(),
				"ip":      c.IP(),
				"ua":      c.Get("User-Agent"),
				"body":    string(c.Body()),
				"query":   c.Queries(),
			}).Info("Inbound request")

			return err
		})
	}

	app.Get("/api/ping", func(c *fiber.Ctx) error {
		v, err := mem.VirtualMemory()
		if err != nil {
			return err
		}

		cpuPercent, err := cpu.Percent(time.Millisecond*100, false)
		if err != nil {
			return err
		}

		uptime, err := host.Uptime()
		if err != nil {
			return err
		}

		var cpuUsage float64
		if len(cpuPercent) > 0 {
			cpuUsage = cpuPercent[0]
		}

		return c.JSON(fiber.Map{
			"status":    "pong",
			"timestamp": time.Now().Unix(),
			"system": fiber.Map{
				"cpu":    cpuUsage,                    // cpu loading
				"ram":    v.UsedPercent,               // ram loading
				"ram_gb": v.Used / 1024 / 1024 / 1024, // gb usage
				"uptime": uptime,                      // server uptime
			},
		})
	})

	api := app.Group("/api")

	public := api.Group("/public")
	authHandler.RegisterRoutes(public)

	private := api.Group("/private")
	private.Use(auth.AuthMiddleware(authSrv, mainRepo))
	private.Use(middleware.UpdateLastLogin(mainRepo))
	private.Use(middleware.LoadRank(ranksSrv))

	authHandler.RegisterPrivateRoute(private)

	userHandler.RegisterRoutes(private)
	inviteHandler.RegisterRoutes(private)
	loggerHandler.RegisterRoutes(private)
	ranksHandler.RegisterRoutes(private)
	notifyHandler.RegisterRoutes(private)
	ticketsHandler.RegisterRoutes(private)
	badgesHandler.RegisterRoutes(private)
	changelogHandler.RegisterRoutes(private)
	developerHandler.RegisterRoutes(private)

	return app, nil
}

func setupLogger() {
	f := slog.NewJSONFormatter()
	f.PrettyPrint = true
	f.TimeFormat = "02/01/2006 15:04:05.000"
	slog.SetFormatter(f)
}
