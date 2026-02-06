package admin

import (
	adminservice "github.com/caseapia/goproject-flush/internal/service/admin"
)

type Handler struct {
	adminService *adminservice.AdminService
}

func NewHandler(r *adminservice.AdminService) *Handler {
	return &Handler{
		adminService: r,
	}
}
