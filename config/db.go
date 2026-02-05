package Config

import (
	database "github.com/caseapia/goproject-flush/internal/db"
	adminrepository "github.com/caseapia/goproject-flush/internal/repository/admin"
	loggerrepo "github.com/caseapia/goproject-flush/internal/repository/logger"
	userrepo "github.com/caseapia/goproject-flush/internal/repository/user"
	"github.com/uptrace/bun"
)

func Connect() *bun.DB {
	if err := database.Connect(); err != nil {
		panic("failed to connect to DB: " + err.Error())
	}
	return database.DB
}

func NewUserRepository() *userrepo.UserRepository {
	return userrepo.NewUserRepository(database.DB)
}

func NewLoggerRepository() *loggerrepo.LoggerRepository {
	return loggerrepo.NewLoggerRepository(database.DB)
}

func NewAdminRepository() *adminrepository.AdminRepository {
	return adminrepository.NewAdminRepository(database.DB)
}
