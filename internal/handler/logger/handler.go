package logger

import service "github.com/caseapia/goproject-flush/internal/service/logger"

type Handler struct {
	service *service.LoggerService
}

func NewHandler(s *service.LoggerService) *Handler {
	return &Handler{service: s}
}
