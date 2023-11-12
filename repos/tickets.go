package repos

import (
	"context"
	"database/sql"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	models "test-manager/usecase_models/boiler"
)

type TicketsRepository interface {
	UpdateTickets(ctx context.Context, Ticket models.Ticket) error
	GetHeadTickets(ctx context.Context, projectId []int) (tickets []*models.Ticket, err error)
	GetTicket(ctx context.Context, id int) (ticket []*models.Ticket, err error)
	SaveTickets(ctx context.Context, Ticket models.Ticket) (int, error)
}

type ticketsRepository struct {
	db *sql.DB
}

func NewTicketsRepository(db *sql.DB) TicketsRepository {
	return &ticketsRepository{db: db}
}

func (r *ticketsRepository) SaveTickets(ctx context.Context, Ticket models.Ticket) (int, error) {
	err := Ticket.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return Ticket.ID, nil
}

func (r *ticketsRepository) UpdateTickets(ctx context.Context, Ticket models.Ticket) error {
	_, err := Ticket.Update(ctx, r.db, boil.Blacklist("account_id", "reply_to", "project_id"))
	if err != nil {
		return err
	}
	return nil
}

func (r *ticketsRepository) GetHeadTickets(ctx context.Context, projectId []int) (tickets []*models.Ticket, err error) {
	var p []interface{}
	for _, value := range projectId {
		p = append(p, value)
	}
	tickets, err = models.Tickets(qm.WhereIn("project_id in ?", p...), models.TicketWhere.ReplyTo.EQ(null.IntFromPtr(nil))).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketsRepository) GetTicket(ctx context.Context, id int) (ticket []*models.Ticket, err error) {
	query := queries.Raw(
		`with recursive cte (id,
									account_id,
									project_id,
									message,
									ticket_status,
									title,
									reply_to,
									created_at,
									updated_at,
									deleted_at) as (
								select id,
									account_id,
									project_id,
									message,
									ticket_status,
									title,
									reply_to,
									created_at,
									updated_at,
									deleted_at
									from tickets
									where id = $1
									and reply_to is null
									union all
									select c.id,
										c.account_id,
										c.project_id,
										c.message,
										c.ticket_status,
										c.title,
										c.reply_to,
										c.created_at,
										c.updated_at,
										c.deleted_at
										from tickets c
										inner join cte
										on c.reply_to = cte.id)
										select *
											from cte;`, id)
	err = query.Bind(ctx, r.db, &ticket)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}
