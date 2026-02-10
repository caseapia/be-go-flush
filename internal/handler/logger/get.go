package logger

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func (l *Handler) GetLogs(c *fiber.Ctx) error {
	logType := c.Params("type")

	var logs interface{}
	var err error

	switch logType {
	case "all":
		logs, err = l.service.GetAllLogs(c.UserContext())
	case "common":
		logs, err = l.service.GetCommonLogs(c.UserContext())
	case "punish":
		logs, err = l.service.GetPunishmentLogs(c.UserContext())
	default:
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "invalid log type, must be 'all', 'common' or 'punish'",
		})
	}

	if err != nil {
		log.Println("Error getting logs:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch logs",
		})
	}

	return c.JSON(logs)
}
