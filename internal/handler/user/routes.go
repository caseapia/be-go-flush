package userhandler

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router, h *UserHandler) {
	users := app.Group("/users")
	user := app.Group("/user")

	users.Get("/", h.GetUsersList)
	user.Get("/:id", h.GetUser)
	user.Put("/admin/:id/ban", h.BanUser)
	user.Delete("/admin/:id/unban", h.UnbanUser)
	user.Put("/admin/create/", h.CreateUser)
	user.Delete("/admin/:id/delete", h.DeleteUser)
	user.Post("/admin/:id/restore", h.RestoreUser)
	user.Patch("/admin/:id/setStatus", h.SetUserStatus)
}
