package logger

import (
	"log"

	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gofiber/fiber/v2"

)

type Handler struct {
	service *logger.Service
}

func NewHandler(s *logger.Service) *Handler {
	return &Handler{service: s}
}

func (l *Handler) SearchLogs(c *fiber.Ctx) error {
	logType := c.Params("type")

	var logs interface{}
	var err error

	switch logType {
	case "common":
		logs, err = l.service.GetCommonLogs(c.UserContext())
	case "punish":
		logs, err = l.service.GetPunishmentLogs(c.UserContext())
	default:
		return &fiber.Error{Code: 404, Message: "invalid log type, must be 'common' or 'punish'"}
	}

	if err != nil {
		log.Println("Error getting logs:", err)
		return &fiber.Error{Code: 500, Message: "failed to fetch logs"}
	}

	return c.JSON(logs)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/admin/logs")

	group.Get("/:type", h.SearchLogs)
}
