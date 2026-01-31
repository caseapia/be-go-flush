package pkgerror

import errormodel "github.com/caseapia/goproject-flush/internal/models/error"

type AppError struct {
	Code    int
	Message errormodel.ErrorMessage
}

func (e *AppError) Error() string {
	return string(e.Message)
}
