package ranks

import (
	"context"
	"fmt"
	"strings"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/pkg/utils/models/enums"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repo   mysql.Repository
	logger logger.Service
}

func NewService(r mysql.Repository, l logger.Service) *Service {
	return &Service{
		repo:   r,
		logger: l,
	}
}

func (s *Service) CreateRank(ctx *fiber.Ctx, adminID uint64, rankName string, rankColor string, rankFlags []string) (*models.Rank, error) {
	u, err := s.repo.SearchUserByID(ctx.UserContext(), adminID)
	if err != nil {
		return nil, err
	}

	r, err := s.repo.SearchRankByID(ctx.UserContext(), int(u.StaffRank))
	if err != nil {
		return nil, err
	}

	if !r.HasFlag("STAFFMANAGEMENT") {
		return nil, fiber.NewError(fiber.StatusForbidden, "you're not allowed to use this function")
	}

	existing, err := s.repo.SearchRankByName(ctx.UserContext(), rankName)
	if existing != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "rank with that name already exists")
	}

	if rankName == "" || len(rankName) < 3 || len(rankName) > 30 {
		return nil, &fiber.Error{Code: 400, Message: "invalid length of rank name"}
	}

	rank := &models.Rank{Name: rankName, Color: rankColor, Flags: rankFlags}

	if err := s.repo.CreateRank(ctx.UserContext(), s.repo.DB, rank); err != nil {
		return nil, err
	}

	addInfo := "with name: " + rankName + ", with color: " + rankColor + "with flags: " + strings.Join(rankFlags, ", ")
	s.logger.Log(ctx.UserContext(), enums.StaffCommonLogger, &adminID, nil, enums.CreateRank, addInfo)

	return rank, nil
}

func (s *Service) DeleteRank(ctx *fiber.Ctx, adminID uint64, id int) (bool, error) {
	u, err := s.repo.SearchUserByID(ctx.UserContext(), adminID)
	if err != nil {
		return false, err
	}
	uRank, err := s.repo.SearchRankByID(ctx.UserContext(), u.StaffRank)
	if err != nil {
		return false, err
	}

	if !uRank.HasFlag("STAFFMANAGEMENT") {
		return false, fiber.NewError(fiber.StatusForbidden, "you're not allowed to use this function")
	}

	r, err := s.repo.SearchRankByID(ctx.UserContext(), id)
	if err != nil {
		return false, err
	}
	if r == nil {
		return false, fiber.NewError(fiber.StatusNotFound, "rank with that name not found")
	}

	if err := s.repo.DeleteRank(ctx.UserContext(), s.repo.DB, r); err != nil {
		return false, err
	}

	addInfo := fmt.Sprintf("with ID: %d, with name: %s", r.ID, r.Name)
	s.logger.Log(ctx.UserContext(), enums.StaffCommonLogger, &adminID, nil, enums.DeleteRank, addInfo)

	return true, nil
}

func (s *Service) SearchAllRanks(ctx *fiber.Ctx) ([]models.Rank, error) {
	ranks, err := s.repo.SearchAllRanks(ctx.UserContext())
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ranks, nil
}

func (s *Service) SearchRankByID(ctx *fiber.Ctx, id int) (*models.Rank, error) {
	rank, err := s.repo.SearchRankByID(ctx.UserContext(), id)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return rank, nil
}

func (s *Service) EditRank(ctx context.Context, sender uint64, rank *models.Rank) (*models.Rank, error) {
	oldRank, err := s.repo.SearchRankByID(ctx, int(rank.ID))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "rank not found")
	}

	oldFlags := strings.Join(oldRank.Flags, ", ")
	oldInfo := fmt.Sprintf("Name: %s, Color: %s, Flags: %v", oldRank.Name, oldRank.Color, oldFlags)

	searchedRank, err := s.repo.SearchRankByName(ctx, rank.Name)
	if err != nil {
		return nil, err
	}
	if searchedRank != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "rank with that name already exists")
	}

	updatedRank, err := s.repo.EditRank(ctx, s.repo.DB, rank)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newFlags := strings.Join(updatedRank.Flags, ", ")
	newInfo := fmt.Sprintf("Name: %s, Color: %s, Flags: %v", updatedRank.Name, updatedRank.Color, newFlags)

	addInfo := "Before: " + oldInfo + "\nAfter: " + newInfo
	s.logger.Log(ctx, enums.StaffCommonLogger, &sender, nil, enums.EditRank, addInfo)

	return updatedRank, nil
}
