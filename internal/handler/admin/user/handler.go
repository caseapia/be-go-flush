package adminUser

import (
	AdminUserService "github.com/caseapia/goproject-flush/internal/service/admin/user"
)

type Handler struct {
	service *AdminUserService.AdminUserService
}

func NewAdminUserHandler(service *AdminUserService.AdminUserService) *Handler {
	return &Handler{service: service}
}
