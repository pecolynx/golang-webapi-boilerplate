package service

import (
	"context"

	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
)

type TicketCreator interface {
	domain.AppUserModel
	AddTicket(ctx context.Context, param TicketAddParameter) (domain.TicketID, error)
	FindMyTickets(ctx context.Context, param TicketSearchCondition) (TicketSearchResult, error)
}

type ticketCreator struct {
	domain.AppUserModel
	rf RepositoryFactory
}

func NewTicketCreator(appUserModel domain.AppUserModel, rf RepositoryFactory) (TicketCreator, error) {
	return &ticketCreator{
		AppUserModel: appUserModel,
		rf:           rf,
	}, nil
}

func (m *ticketCreator) AddTicket(ctx context.Context, param TicketAddParameter) (domain.TicketID, error) {
	ticketRepo := m.rf.NewTicketRepository(ctx)
	ticketID, err := ticketRepo.AddTicket(ctx, m.GetAppUserID(), param)
	if err != nil {
		return nil, err
	}

	return ticketID, nil
}

func (m *ticketCreator) FindMyTickets(ctx context.Context, param TicketSearchCondition) (TicketSearchResult, error) {
	ticketRepo := m.rf.NewTicketRepository(ctx)
	result, err := ticketRepo.FindMyTickets(ctx, m.GetAppUserID(), param)
	if err != nil {
		return nil, err
	}

	return result, nil
}
