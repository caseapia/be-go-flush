package models

import (
	"time"

	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	"github.com/uptrace/bun"
)

type Ticket struct {
	bun.BaseModel `bun:"table:tickets"`

	ID        uint64               `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time            `bun:"created_at,notnull" json:"created_at"`
	Status    enums.TicketStatus   `bun:"status,notnull" json:"status"`
	Priority  enums.TicketPriority `bun:"priority" json:"prioriy"`
	AuthorID  uint64               `bun:"author_id,notnull" json:"author_id"`
	HandledBy *uint64              `bun:"handling_by" json:"-"`
	Title     string               `bun:"title,notnull" json:"title"`
	Category  string               `bun:"category,notnull" json:"category"`
	UpdatedAt time.Time            `bun:"updated_at" json:"updated_at"`

	Author  *TicketAuthor `bun:"rel:belongs-to,join:author_id=id" json:"author"`
	Handler *UserRelation `bun:"rel:belongs-to,join:handling_by=id" json:"handler"`
}
type TicketMessage struct {
	bun.BaseModel `bun:"table:tickets_messages"`

	ID        uint64    `bun:"id,pk,autoincrement" json:"id"`
	TicketID  uint64    `bun:"ticket_id,notnull" json:"ticket_id"`
	AuthorID  uint64    `bun:"author_id,notnull" json:"-"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"created_at"`
	Content   string    `bun:"content,notnull" json:"content"`

	Author *UserRelation `bun:"rel:belongs-to,join:author_id=id" json:"author"`
}

type TicketAction struct {
	bun.BaseModel `bun:"table:flush_db_logs.tickets_actions"`

	ID        uint64    `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	TicketID  uint64    `bun:"ticket_id" json:"ticket_id"`
	AuthorID  uint64    `bun:"author_id" json:"-"`
	Action    string    `bun:"action" json:"action"`

	Author *UserRelation `bun:"rel:belongs-to,join:author_id=id" json:"author"`
}
type TicketCreationInput struct {
	Title        string `bun:"title,notnull" json:"title"`
	Category     string `bun:"category,notnull" json:"category"`
	FirstMessage string `json:"message"`
}

type TicketCategoryChangingInput struct {
	NewCategory string `bun:"category,notnull" json:"new_category"`
}

type TicketMessageCreationInput struct {
	Ticket   Ticket `json:"ticket"`
	AuthorID uint64 `bun:"author_id,notnull" json:"author_id"`
	Content  string `bun:"content,notnull" json:"content"`
}

type TicketAuthor struct {
	bun.BaseModel `bun:"table:flush_db.users"`

	UserRelation
	LastLogin *time.Time `bun:"last_login" json:"last_login"`
}

type TicketResponse struct {
	Ticket   Ticket          `json:"ticket"`
	Messages []TicketMessage `json:"messages"`
	Action   []TicketAction  `json:"actions"`
}
