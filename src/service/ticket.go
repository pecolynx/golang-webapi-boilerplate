package service

import (
	"context"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"

	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
)

type Ticket interface {
	domain.TicketModel
}

type ticket struct {
	domain.TicketModel
}

func NewTicket(ctx context.Context, ticketModel domain.TicketModel) (Ticket, error) {
	m := &ticket{
		TicketModel: ticketModel,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *ticket) GetTicketModel() domain.TicketModel {
	return m.TicketModel
}
