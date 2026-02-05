package AdminUserHandler

import (
	AdminUserService "github.com/caseapia/goproject-flush/internal/service/admin/user"
)

type AdminUserHandler struct {
	service *AdminUserService.AdminUserService
}

func NewAdminUserHandler(service *AdminUserService.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{service: service}
}
