package pkgerror

import (
	errormodel "github.com/caseapia/goproject-flush/internal/models/error"
	"github.com/gookit/slog"
)

type AppError struct {
	Code    int
	Message errormodel.ErrorMessage
}

func (e *AppError) Error() string {
	slog.WithData(slog.M{
		"properties": e,
	}).Error("AppError encountered")

	return string(e.Message)
}
