package tickets

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/tickets")
	groupAdmin := router.Group("/admin/tickets")

	group.Post("/create", h.CreateTicket)                                                     // & Create ticket
	group.Get("/populate/:id", h.PopulateTicket)                                              // & Populate one ticket (only staff members can populate tickets created by other user)
	group.Get("/mytickets", h.PopulateAllUserTickets)                                         // & Populate all tickets created by selected user
	group.Post("/send", h.CreateTicketMessage)                                                // & Create message in a ticket
	group.Patch("/close/:id", h.CloseTicket)                                                  // & Close ticket
	groupAdmin.Get("/populate", middleware.RequireFlag("STAFF"), h.SearchTickets)             // ~ Populate all tickets existed in database
	groupAdmin.Post("/assign/:id", middleware.RequireFlag("STAFF"), h.TicketAssignment)       // ~ Assign an admin to the ticket
	groupAdmin.Patch("/edit/category/:id", middleware.RequireFlag("STAFF"), h.ChangeCategory) // ~ Change category of ticket
}
