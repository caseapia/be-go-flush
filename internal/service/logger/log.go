package LoggerService

import (
	"context"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
)

func (l *LoggerService) Log(
	ctx context.Context,
	adminID uint64,
	userID *uint64,
	action loggermodule.LoggerAction,
	additional ...string,
) error {
	var addInfo *string
	if len(additional) > 0 {
		addInfo = &additional[0]
	}

	var u *usermodel.User
	if userID != nil {
		var err error
		u, err = l.uRepo.GetByID(ctx, *userID)
		if err != nil {
			u = nil
		}
	}

	logEntry := loggermodule.ActionLog{
		AdminID:        adminID,
		Action:         action,
		AdditionalInfo: addInfo,
	}

	if userID != nil && u != nil {
		logEntry.UserID = userID
		logEntry.UserNickname = &u.Name
	}

	return l.repo.Log(ctx, &logEntry)
}
