package repository

import (
	"fmt"
	"todo/internal/model"
	"todo/pkg/db"
)

type TicketRepository struct {
	db *db.MySQL
}

func NewTicketRepository(db *db.MySQL) *TicketRepository {
	return &TicketRepository{db}
}

func (r *TicketRepository) GetAll() ([]model.Ticket, error) {
	rows, err := r.db.Select(fmt.Sprintf("select * from tickets where status <> %d", model.TICKET_STATUS_DELITED))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Ticket
	for rows.Next() {
		var ticket model.Ticket
		err = rows.Scan(&ticket.ID, &ticket.Title, &ticket.Status, &ticket.CreatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, ticket)
	}

	return result, nil
}

func (r *TicketRepository) GetById(id int) (*model.Ticket, error) {
	rows, err := r.db.Select(fmt.Sprintf("select * from tickets where id = %d", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ticket model.Ticket
		err = rows.Scan(&ticket.ID, &ticket.Title, &ticket.Status, &ticket.CreatedAt)
		if err != nil {
			return nil, err
		}

		return &ticket, nil
	}

	return nil, nil
}

func (r *TicketRepository) Create(ticket *model.Ticket) error {
	id, err := r.db.Insert(fmt.Sprintf("insert into tickets (title, status) values ('%s', %d)", ticket.Title, ticket.Status))
	if err != nil {
		return err
	}
	ticket.ID = id

	return nil
}

func (r *TicketRepository) Update(ticket *model.Ticket) error {
	_, err := r.db.Update(fmt.Sprintf("update tickets set title = '%s', status = %d where id = %d", ticket.Title, ticket.Status, ticket.ID))
	if err != nil {
		return err
	}

	return nil
}
