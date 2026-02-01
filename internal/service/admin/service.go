package adminservice

import (
	adminrepository "github.com/caseapia/goproject-flush/internal/repository/admin"
)

type AdminService struct {
	adminRepo *adminrepository.AdminRepository
}

func NewAdminService(r *adminrepository.AdminRepository) *AdminService {
	return &AdminService{
		adminRepo: r,
	}
}
