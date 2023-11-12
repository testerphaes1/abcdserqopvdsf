package usecase_models

import (
	"github.com/volatiletech/null/v8"
	"time"
)

const (
	TicketStatusPending  = 0
	TicketStatusResolves = 1
)

type Tickets struct {
	ID           int         `json:"id"`
	AccountID    int         `json:"account_id"`
	ProjectID    null.Int    `json:"project_id"`
	Message      null.String `json:"message"`
	TicketStatus int         `json:"ticket_status"`
	Title        null.String `json:"title"`
	ReplyTo      null.Int    `json:"reply_to"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	DeletedAt    null.Time   `json:"deleted_at"`
}

type CreateTicketsResponse struct {
	TicketId int `json:"ticket_id"`
}
