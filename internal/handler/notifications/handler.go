package notifications

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/notifications"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *notifications.Service
}

func NewHandler(s *notifications.Service) *Handler {
	return &Handler{service: s}
}

func (s *Handler) PopulateNotifications(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	notifications, err := s.service.PopulateNotifications(c.UserContext(), sender.ID, sender.ID)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, notifications)
}

func (s *Handler) ReadNotifications(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	notifications, err := s.service.ReadNotifications(c.Context(), sender.ID)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, notifications)
}

func (s *Handler) ClearNotifications(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	notifications, _ := s.service.ClearNotifications(c.Context(), sender.ID)

	return utils.Success(c, 200, notifications)
}

func (s *Handler) RemoveOwnNotification(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	var input models.RemoveNotificationRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	isDeleted, err := s.service.RemoveNotification(c.UserContext(), sender.ID, sender.ID, input.NotifyID)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, isDeleted)
}

// ! Admin actions
func (s *Handler) SendNotification(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	sender := account.GetUserFromContext(c)

	var input models.SendNotificationRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	s.service.SendNotification(c.UserContext(), uint64(id), input.Type, input.Title, input.Text, &sender.ID)

	return utils.Success(c, 200, fiber.Map{"status": "success"})
}

func (s *Handler) PopulateUserNotifications(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	sender := account.GetUserFromContext(c)

	notifications, err := s.service.PopulateNotifications(c.UserContext(), uint64(id), sender.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, 200, notifications)
}

func (s *Handler) RemoveNotification(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	sender := account.GetUserFromContext(c)

	var input models.RemoveNotificationRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	isDeleted, err := s.service.RemoveNotification(c.UserContext(), uint64(id), sender.ID, input.NotifyID)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, isDeleted)
}
