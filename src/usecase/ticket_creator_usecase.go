package usecase

import (
	"context"

	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
)

const DefaultPageNo = 1
const DefaultPageSize = 10

type TicketCreatorUsecase interface {
	AddTicket(ctx context.Context, operatorID domain.TicketCreatorID, parameter service.TicketAddParameter) (domain.TicketID, error)
}

type ticketCreatorUsecase struct {
	transactionManager service.TransactionManager
}

func NewTicketCreatorUsecase(transactionManager service.TransactionManager) TicketCreatorUsecase {
	return &ticketCreatorUsecase{
		transactionManager: transactionManager,
	}
}

func (s *ticketCreatorUsecase) AddTicket(ctx context.Context, operatorID domain.TicketCreatorID, parameter service.TicketAddParameter) (domain.TicketID, error) {
	var addedTicketID domain.TicketID

	if err := s.transactionManager.Do(ctx, func(rf service.RepositoryFactory) error {
		appUserRepo := rf.NewAppUserRepository(ctx)
		ticketCreator, err := appUserRepo.FindTicketCreatorByID(ctx, operatorID)
		if err != nil {
			return err
		}
		tmpAddedTicketID, err := ticketCreator.AddTicket(ctx, parameter)
		if err != nil {
			return err
		}
		addedTicketID = tmpAddedTicketID
		return nil
	}); err != nil {
		return nil, err
	}

	return addedTicketID, nil
}
