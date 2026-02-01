package adminhandler

import (
	adminservice "github.com/caseapia/goproject-flush/internal/service/admin"
)

type AdminHandler struct {
	adminService *adminservice.AdminService
}

func NewAdminHandler(r *adminservice.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: r,
	}
}
