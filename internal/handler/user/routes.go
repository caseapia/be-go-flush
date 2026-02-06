package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router, h *UserHandler) {
	users := app.Group("/users")
	user := app.Group("/user")

	users.Get("/", h.GetUsersList)
	user.Get("/:id", h.GetUser)
}
