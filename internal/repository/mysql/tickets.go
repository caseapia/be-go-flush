package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (r *Repository) SearchOpenedTickets(ctx context.Context) ([]models.Ticket, int, error) {
	var tickets []models.Ticket

	query := r.DB.NewSelect().
		Model(&tickets).
		Where("status != ?", "closed").
		Order("created_at ASC").
		Relation("Author").
		Relation("Handler").
		Limit(COLUMNS_LIMIT)

	err := query.Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return tickets, COLUMNS_LIMIT, nil
}

func (r *Repository) SearchAllTickets(ctx context.Context) ([]models.Ticket, int, error) {
	var tickets []models.Ticket

	query := r.DB.NewSelect().
		Model(&tickets).
		Order("created_at ASC").
		Relation("Author").
		Relation("Handler").
		Limit(COLUMNS_LIMIT)

	err := query.Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return tickets, COLUMNS_LIMIT, nil
}

func (r *Repository) SearchTicketByID(ctx context.Context, ticketID uint64) (*models.Ticket, *[]models.TicketMessage, error) {
	var ticket models.Ticket
	var messages []models.TicketMessage

	err := r.DB.NewSelect().
		Model(&ticket).
		Where("ticket.id = ?", ticketID).
		Relation("Author").
		Relation("Handler").
		Scan(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = r.DB.NewSelect().
		Model(&messages).
		Where("ticket_id = ?", ticketID).
		Relation("Author").
		Scan(ctx)
	if err != nil {
		return nil, nil, err
	}

	if messages == nil {
		messages = make([]models.TicketMessage, 0)
	}

	return &ticket, &messages, nil
}

func (r *Repository) PopulateTicket(ctx context.Context, ticketID uint64) (*models.Ticket, error) {
	t := new(models.Ticket)

	err := r.DB.NewSelect().
		Model(t).
		Relation("Author").
		Relation("Handler").
		Where("ticket.id = ?", ticketID).
		Limit(1).
		Scan(ctx)

	return t, err
}

func (r *Repository) TicketAssignment(ctx context.Context, tx bun.IDB, ticketID uint64, userID uint64) error {
	_, err := tx.NewUpdate().
		Model((*models.Ticket)(nil)).
		Set("handling_by = ?", userID).
		Where("id = ?", ticketID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) TicketUnassignment(ctx context.Context, tx bun.IDB, ticketID uint64) error {
	_, err := tx.NewUpdate().
		Model((*models.Ticket)(nil)).
		Set("handling_by = ?", nil).
		Where("id = ?", ticketID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) PopulateAllUserTickets(ctx context.Context, userID uint64) ([]models.Ticket, error) {
	var tickets []models.Ticket

	err := r.DB.NewSelect().
		Model(&tickets).
		Relation("Author").
		Relation("Handler").
		Where("author_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	if tickets == nil {
		tickets = make([]models.Ticket, 0)
	}

	return tickets, nil
}

func (r *Repository) CreateTicket(ctx context.Context, tx bun.IDB, entry models.Ticket) (*models.Ticket, error) {
	_, err := tx.NewInsert().
		Model(&entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to create ticket!")

		return nil, err
	}

	return &entry, nil
}

func (r *Repository) SetTicketStatus(ctx context.Context, tx bun.IDB, ticketID uint64, newStatus enums.TicketStatus) (*models.Ticket, error) {
	t := new(models.Ticket)

	_, err := tx.NewUpdate().
		Model(t).
		Set("status = ?", newStatus).
		Where("id = ?", ticketID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *Repository) PopulateTicketMessages(ctx context.Context, ticketID uint64) ([]models.TicketMessage, error) {
	var messages []models.TicketMessage

	query := r.DB.NewSelect().
		Model(&messages).
		Relation("Author.Staff").
		Relation("Author.Developer").
		Where("ticket_id = ?", ticketID)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	if messages == nil {
		messages = make([]models.TicketMessage, 0)
	}

	return messages, err
}

func (r *Repository) CreateTicketMessage(ctx context.Context, tx bun.IDB, entry models.TicketMessage, handlerID *uint64) error {
	_, err := tx.NewInsert().
		Model(&entry).
		Returning("*").
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err, "entry": entry}).Error("failed to send message in ticket")

		return err
	}

	var newStatus enums.TicketStatus
	if handlerID != nil && entry.AuthorID == *handlerID {
		newStatus = enums.TicketStatusOpen
	} else {
		newStatus = enums.TicketStatusPending
	}

	_, statusErr := r.SetTicketStatus(ctx, tx, uint64(entry.TicketID), newStatus)
	if statusErr != nil {
		return statusErr
	}

	return nil
}

func (r *Repository) CloseTicket(ctx context.Context, tx bun.IDB, ticketID uint64) error {
	_, err := tx.NewUpdate().
		Model((*models.Ticket)(nil)).
		Set("status = ?", enums.TicketStatusClosed).
		Where("id = ?", ticketID).
		Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ChangeTicketCategory(ctx context.Context, tx bun.IDB, ticketID uint64, newCategory string) error {
	_, err := r.DB.NewUpdate().
		Model((*models.Ticket)(nil)).
		Set("category = ?", newCategory).
		Where("id = ?", ticketID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteTicket(ctx context.Context, tx bun.IDB, ticketID uint64) error {
	_, err := tx.NewDelete().
		Model((*models.Ticket)(nil)).
		Where("id = ?", ticketID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) PopulateTicketActions(ctx context.Context, ticketID uint64) ([]models.TicketAction, error) {
	var actions []models.TicketAction

	err := r.DB.NewSelect().
		Model(&actions).
		Where("ticket_id = ?", ticketID).
		Relation("Author.Staff").
		Relation("Author.Developer").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	if actions == nil {
		actions = make([]models.TicketAction, 0)
	}

	return actions, nil
}

func (r *Repository) CreateTicketAction(ctx context.Context, tx bun.IDB, entry models.TicketAction) ([]models.TicketAction, error) {
	_, err := tx.NewInsert().
		Model(&entry).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	actions, err := r.PopulateTicketActions(ctx, entry.TicketID)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
