package tickets

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/internal/service/user/notifications"
	"github.com/caseapia/goproject-flush/pkg/utils/models/enums"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

type Service struct {
	repo   mysql.Repository
	notify notifications.Service
	logger logger.Service
}

func NewService(r mysql.Repository, n notifications.Service, l logger.Service) *Service {
	return &Service{
		repo:   r,
		notify: n,
		logger: l,
	}
}

func (s *Service) createTicketAction(ctx context.Context, ticketID, authorID uint64, action string) ([]models.TicketAction, error) {
	actions, err := s.repo.CreateTicketAction(ctx, s.repo.DB, models.TicketAction{
		CreatedAt: time.Now(),
		TicketID:  ticketID,
		AuthorID:  authorID,
		Action:    action,
	})
	if err != nil {
		return nil, err
	}

	return actions, nil
}

func (s *Service) SearchTickets(ctx context.Context, user *models.User) ([]models.Ticket, int, error) {
	staffRank, err := s.repo.SearchRankByID(ctx, user.StaffRank)
	if err != nil {
		return nil, 0, err
	}

	developerRank, err := s.repo.SearchRankByID(ctx, user.DeveloperRank)
	if err != nil {
		return nil, 0, err
	}

	isUserHasPermissionToViewClosedTickets := user.UserHasFlag("SENIOR")
	isRankHasPermissionToViewClosedTickets := staffRank.HasFlag("SENIOR") || developerRank.HasFlag("LEADDEV")

	var tickets []models.Ticket
	var columns int

	if !isUserHasPermissionToViewClosedTickets && !isRankHasPermissionToViewClosedTickets {
		openedTickets, openedTicketsColumns, err := s.repo.SearchOpenedTickets(ctx)
		if err != nil {
			return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tickets = openedTickets
		columns = openedTicketsColumns
	} else {
		closedTickets, closedTicketsColumns, err := s.repo.SearchAllTickets(ctx)
		if err != nil {
			return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tickets = closedTickets
		columns = closedTicketsColumns
	}

	return tickets, columns, nil
}

func (s *Service) PopulateTicket(ctx context.Context, ticketID uint64, user *models.User) (*models.TicketResponse, error) {
	ticket, err := s.repo.PopulateTicket(ctx, ticketID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	staffRank, err := s.repo.SearchRankByID(ctx, user.StaffRank)
	if err != nil {
		return nil, err
	}

	developerRank, err := s.repo.SearchRankByID(ctx, user.DeveloperRank)
	if err != nil {
		return nil, err
	}

	isAuthor := ticket.Author.ID == user.ID
	isStaff := user.UserHasFlag("STAFF")
	rankIsStaff := staffRank.HasFlag("STAFF") || developerRank.HasFlag("STAFF")
	isUserHasPermissionToViewClosedTickets := user.UserHasFlag("SENIOR")
	isRankHasPermissionToViewClosedTickets := staffRank.HasFlag("SENIOR") || developerRank.HasFlag("LEADDEV")

	if !isAuthor && !(isStaff || rankIsStaff) {
		return nil, fiber.NewError(fiber.StatusForbidden, "you have no access to this ticket")
	}

	if !isAuthor && ticket.Status == enums.TicketStatusClosed && !(isUserHasPermissionToViewClosedTickets || isRankHasPermissionToViewClosedTickets) {
		return nil, fiber.NewError(fiber.StatusNotFound, "selected ticket was not found")
	}

	messages, err := s.repo.PopulateTicketMessages(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	actions, err := s.repo.PopulateTicketActions(ctx, ticketID)

	return &models.TicketResponse{
		Ticket:   *ticket,
		Messages: messages,
		Action:   actions,
	}, nil
}

func (s *Service) PopulateAllUserTickets(ctx context.Context, userID uint64) ([]models.Ticket, error) {
	tickets, err := s.repo.PopulateAllUserTickets(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return tickets, nil
}

func (s *Service) CreateTicket(ctx context.Context, user models.User, title, category, message string) (*models.Ticket, error) {
	// parameters
	categories := []string{
		"Technical Support",
		"Billing & Payments",
		"Account Access",
		"Report a Bug",
		"Other",
	}

	// terms
	if len(title) > 255 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "title too long")
	}
	if !slices.Contains(categories, category) {
		return nil, fiber.NewError(fiber.StatusNotFound, "category not found")
	}

	var ticket *models.Ticket
	var err error
	s.repo.WithTx(ctx, func(tx bun.Tx) error {
		// Ticket creation
		ticket, err = s.repo.CreateTicket(ctx, tx, models.Ticket{
			AuthorID:  user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     title,
			Category:  category,
			HandledBy: nil,
			Status:    enums.TicketStatusPending,
			Priority:  enums.TicketPriorityLow,
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		var handlerID *uint64
		if ticket.Handler != nil {
			handlerID = &ticket.Handler.ID
		}

		// Message creation
		if ticket != nil {
			mErr := s.repo.CreateTicketMessage(ctx, tx, models.TicketMessage{
				TicketID:  ticket.ID,
				AuthorID:  user.ID,
				CreatedAt: time.Now(),
				Content:   message,
			}, handlerID)
			if mErr != nil {
				return fiber.NewError(fiber.StatusInternalServerError, mErr.Error())
			}
		}

		return nil
	})

	return ticket, nil
}

func (s *Service) CreateTicketMessage(ctx context.Context, ticket *models.Ticket, user *models.User, content string) (*models.TicketResponse, error) {
	rank, err := s.repo.SearchRankByID(ctx, user.StaffRank)
	if err != nil {
		return nil, err
	}

	isAuthor := ticket.Author.ID == user.ID
	isStaff := user.UserHasFlag("STAFF")
	rankIsStaff := rank.HasFlag("STAFF")
	isStaffManagement := user.UserHasFlag("STAFFMANAGEMENT")
	rankIsStaffManagement := rank.HasFlag("STAFFMANAGEMENT")
	isHandler := ticket.Handler != nil && ticket.Handler.ID == user.ID

	if !isAuthor && !(isStaff || rankIsStaff) {
		slog.WithData(slog.M{
			"isAuthor":    isAuthor,
			"isStaff":     isStaff,
			"rankIsStaff": rankIsStaff,
			"ticket":      ticket,
			"user.ID":     user.ID,
		}).Error("error when sending message")
		return nil, fiber.NewError(fiber.StatusForbidden, "you have no access to this ticket")
	}
	if ticket.Status == enums.TicketStatusClosed && !(isStaffManagement || rankIsStaffManagement) {
		return nil, fiber.NewError(fiber.StatusNotAcceptable, "the ticket is closed and you can't answer here")
	}

	if !isAuthor && !isHandler {
		slog.WithData(slog.M{
			"isAuthor": isAuthor,
			"ticket":   ticket,
			"user":     user,
		}).Error("message was not sent")

		return nil, fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("this ticket already handled by %s and you cannot answer here", ticket.Handler.Name))
	}

	var updatedTicket *models.TicketResponse

	var handlerID *uint64
	if ticket.Handler != nil {
		handlerID = &ticket.Handler.ID
	}

	s.repo.WithTx(ctx, func(tx bun.Tx) error {
		err = s.repo.CreateTicketMessage(ctx, tx, models.TicketMessage{
			TicketID:  ticket.ID,
			AuthorID:  user.ID,
			CreatedAt: time.Now(),
			Content:   content,
		}, handlerID)
		if err != nil {
			return err
		}

		updatedTicket, err = s.PopulateTicket(ctx, ticket.ID, user)
		if err != nil {
			return err
		}

		return nil
	})

	return updatedTicket, nil
}

func (s *Service) TicketAssignment(ctx context.Context, ticketID uint64, user *models.User) (*models.TicketResponse, error) {
	rank, err := s.repo.SearchRankByID(ctx, user.StaffRank)
	if err != nil {
		return nil, err
	}

	isStaff := user.UserHasFlag("STAFF")
	rankIsStaff := rank.HasFlag("STAFF")

	if !(isStaff || rankIsStaff) {
		return nil, fiber.NewError(fiber.StatusForbidden, "not allowed to use this function")
	}

	ticket, ticketErr := s.PopulateTicket(ctx, ticketID, user)
	if ticketErr != nil {
		return nil, ticketErr
	}

	if ticket.Ticket.HandledBy != nil {
		unassignmentErr := s.repo.TicketUnassignment(ctx, s.repo.DB, ticketID)
		if unassignmentErr != nil {
			return nil, err
		}

		s.notify.SendNotification(ctx, ticket.Ticket.AuthorID, enums.Information, "Your ticket has been updated", fmt.Sprintf("Staff member has unassign himself from your ticket #%v", ticket.Ticket.ID), nil)
		s.logger.Log(ctx, enums.AdminTicketLogger, &user.ID, &ticket.Ticket.Author.ID, enums.AssignedToTicket, fmt.Sprintf("ID: %v | Title: %v", ticketID, ticket.Ticket.Title))
		s.createTicketAction(ctx, ticketID, user.ID, string(enums.TicketActionUnassign))
	} else {
		assignmentErr := s.repo.TicketAssignment(ctx, s.repo.DB, ticketID, user.ID)
		if assignmentErr != nil {
			return nil, assignmentErr
		}

		s.notify.SendNotification(ctx, ticket.Ticket.Author.ID, enums.Success, "Your ticket has been updated", fmt.Sprintf("You have a new staff member, assigned for your ticket #%v", ticket.Ticket.ID), nil)
		s.logger.Log(ctx, enums.AdminTicketLogger, &user.ID, &ticket.Ticket.Author.ID, enums.UnassignFromTicket, fmt.Sprintf("ID: %v | Title: %v", ticketID, ticket.Ticket.Title))
		s.createTicketAction(ctx, ticketID, user.ID, string(enums.TicketActionHandled))
	}

	return ticket, err
}

func (s *Service) CloseTicket(ctx context.Context, ticketID uint64, user models.User) (*models.TicketResponse, error) {
	rank, err := s.repo.SearchRankByID(ctx, user.StaffRank)
	if err != nil {
		return nil, err
	}

	ticket, _, err := s.repo.SearchTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	isAuthor := ticket.Author.ID == user.ID
	isStaffManagement := user.UserHasFlag("STAFFMANAGEMENT")
	rankIsStaffManagement := rank.HasFlag("STAFFMANAGEMENT")
	isHandler := ticket.Handler != nil && ticket.Handler.ID == user.ID

	if !((isAuthor || isHandler) || isStaffManagement || rankIsStaffManagement) {
		return nil, fiber.NewError(fiber.StatusForbidden, "you are not allowed to use this action")
	}
	if ticket.Status == enums.TicketStatusClosed {
		return nil, fiber.NewError(fiber.StatusBadRequest, "ticket already closed")
	}

	err = s.repo.CloseTicket(ctx, s.repo.DB, ticket.ID)
	if err != nil {
		return nil, err
	}

	if !isAuthor {
		s.notify.SendNotification(ctx, ticket.Author.ID, enums.Success, "Your ticket has been updated", fmt.Sprintf("Your ticket #%v was closed by an admin", ticket.ID), nil)
		s.logger.Log(ctx, enums.AdminTicketLogger, &user.ID, &ticket.Author.ID, enums.CloseTicket, fmt.Sprintf("ID: %v | Title: %v", ticketID, ticket.Title))
	}

	s.createTicketAction(ctx, ticketID, user.ID, string(enums.TicketActionClose))

	updatedTicket, err := s.PopulateTicket(ctx, ticket.ID, &user)
	if err != nil {
		return nil, err
	}

	return updatedTicket, nil
}

func (s *Service) ChangeTicketCategory(ctx context.Context, ticketID uint64, newCategory string, user *models.User) (*models.TicketResponse, error) {
	ticket, err := s.repo.PopulateTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	isHandler := ticket.Handler != nil && ticket.Handler.ID == user.ID
	if !isHandler {
		return nil, fiber.NewError(fiber.StatusForbidden, "category of this ticket can only be changed by the staff member who is handling this ticket")
	}

	oldCategory := ticket.Category

	changeCategoryErr := s.repo.ChangeTicketCategory(ctx, s.repo.DB, ticket.ID, newCategory)
	if changeCategoryErr != nil {
		return nil, changeCategoryErr
	}

	ticket.Category = newCategory

	actionMsg := fmt.Sprintf(string(enums.TicketActionChangeCategory), newCategory)

	s.notify.SendNotification(ctx, ticket.AuthorID, enums.Information, "Your ticket has been updated", fmt.Sprintf("Category of your ticket #%v has been changed by an admin to %s", ticket.ID, ticket.Category), nil)
	s.logger.Log(ctx, enums.AdminTicketLogger, &user.ID, &ticket.Author.ID, enums.ChangeTicketCategory, fmt.Sprintf("ID: %v | Title: %s\nCategory before: %s | Category after: %s", ticketID, ticket.Title, oldCategory, newCategory))
	s.createTicketAction(ctx, ticketID, user.ID, actionMsg)

	updatedTicket, err := s.PopulateTicket(ctx, ticket.ID, user)
	if err != nil {
		return nil, err
	}

	return updatedTicket, nil
}
