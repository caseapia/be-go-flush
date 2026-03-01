package auth

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/auth") // & Core route

	group.Post("/refresh", h.Refresh)   // & Refresh access token
	group.Post("/register", h.Register) // & Register account
	group.Post("/login", h.Login)       // & Login in existing account
}

func (h *Handler) RegisterPrivateRoute(router fiber.Router) {
	group := router.Group("/auth") // & Core route

	group.Get("/discord", h.DiscordRedirect)           // & Get discord OTP Link
	group.Post("/discord/callback", h.DiscordCallback) // & Link discord
	group.Delete("/discord/unlink", h.DiscordUnlink)   // & Unlink discord
	group.Delete("/logout", h.Logout)                  // & Logout
}
