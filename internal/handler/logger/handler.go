package logger

import (
	"log"

	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/models"
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
	var logs interface{}
	var limit int
	var err error

	uVal := c.Locals("user")
	_, ok := uVal.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.LogPopulate

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	keywords := ""
	if input.Keywords != nil {
		keywords = *input.Keywords
	}

	switch input.Type {
	case "staffCommon":
		logs, limit, err = l.service.GetCommonStaffLogs(c.UserContext(), input.StartDate, input.EndDate, keywords)
	case "staffPunish":
		logs, limit, err = l.service.GetPunishmentStaffLogs(c.UserContext(), input.StartDate, input.EndDate, keywords)
	case "adminTickets":
		logs, limit, err = l.service.GetTicketsAdminLogs(c.UserContext(), input.StartDate, input.EndDate, keywords)
	case "adminAuth":
		logs, limit, err = l.service.GetAuthAdminLogs(c.UserContext(), input.StartDate, input.EndDate, keywords)
	default:
		return fiber.NewError(fiber.StatusNotFound, "invalid log type")
	}

	logs = struct {
		Data  interface{} `json:"data"`
		Limit int         `json:"limit"`
	}{
		Data:  logs,
		Limit: limit,
	}

	if err != nil {
		log.Println("Error getting logs:", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(logs)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	adminGroup := router.Group("/admin/logs")
	staffGroup := router.Group("/admin/logs/staff")

	adminGroup.Post("/populate", middleware.RequireFlag("ADMIN"), h.SearchLogs)
	staffGroup.Post("/populate", middleware.RequireFlag("STAFFMANAGEMENT"), h.SearchLogs)
}
