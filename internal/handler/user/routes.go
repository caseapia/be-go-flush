package user

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/user")
	groupAdmin := router.Group("/admin/user")
	groupSupport := router.Group("/support/user")

	group.Get("/account", h.GetOwnAccount)                                    // & Get information about authorized account
	groupAdmin.Get("/all", middleware.RequireFlag("ADMIN"), h.SearchAllUsers) // ~ Search by all accounts
	group.Get("/:id", h.SearchUserByID)                                       // & Get information about account by ID (Only staff members can bypass privacy settings)

	groupAdmin.Put("/create", middleware.RequireFlag("SENIOR"), h.CreateUser)                                 // ~ Create a new user
	groupAdmin.Patch("/ban/:id", middleware.RequireFlag("ADMIN"), h.BanUser)                                  // ~ Ban user
	groupAdmin.Delete("/unban/:id", middleware.RequireFlag("ADMIN"), h.UnbanUser)                             // ~ Unban user
	groupAdmin.Delete("/delete/:id", middleware.RequireFlag("SENIOR"), h.DeleteUser)                          // ~ Delete user (soft and hard deletion)
	groupAdmin.Put("/restore/:id", middleware.RequireFlag("MANAGER"), h.RestoreUser)                          // ~ Restore soft deleted user
	groupAdmin.Patch("/rank/staff/:id", middleware.RequireFlag("STAFFMANAGEMENT"), h.SetStaffRank)            // ~ Set staff rank
	groupAdmin.Patch("/rank/developer/:id", middleware.RequireFlag("MANAGER"), h.SetDeveloperRank)            // ~ Set developer rank
	groupAdmin.Get("/:id", middleware.RequireFlag("ADMIN"), h.GetUserPrivate)                                 // ~ Get sensetive user information
	groupAdmin.Patch("/edit/:id", middleware.RequireFlag("MANAGER"), h.ChangeUser)                            // ~ Change user settings (like e-mail, nickname or password)
	groupAdmin.Patch("/editflag/:id", middleware.RequireFlag("STAFFMANAGEMENT"), h.EditUserFlags)             // ~ Edit user access flags (similar with rank flags)
	groupAdmin.Delete("/reset/:id", middleware.RequireFlag("SENIOR"), h.ResetUserSensetiveData)               // ~ Reset user sensetive data
	groupAdmin.Patch("/editbadges/:id", middleware.RequireFlag("SENIOR"), h.EditUserBadges)                   // ~ Edit user badges
	groupAdmin.Delete("/discord/forceunlink/:id", middleware.RequireFlag("SENIOR"), h.ForceUnlinkUserDiscord) // ~ Force break relation of discord and website account
	groupSupport.Get("/bans/populate", middleware.RequireFlag("STAFF"), h.PopulateBanList)                    // ~ Populate all bans ever issued
}
