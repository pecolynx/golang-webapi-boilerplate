package usecase_test

import (
	"context"
	"testing"

	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	domain_mock "github.com/pecolynx/golang-webapi-boilerplate/src/domain/mocks"
	"github.com/pecolynx/golang-webapi-boilerplate/src/gateway"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
	service_mock "github.com/pecolynx/golang-webapi-boilerplate/src/service/mocks"
	"github.com/pecolynx/golang-webapi-boilerplate/src/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ticketCreatorUsecase_AddTicket(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	ticketCreatorID, err := domain.NewTicketCreatorID(123)
	require.NoError(t, err)
	parameter, err := service.NewTicketAddParameter("TITLE", "DESCRIPTION")
	require.NoError(t, err)
	ticketID, err := domain.NewTicketID(456)
	require.NoError(t, err)

	type input struct {
		operatorID domain.TicketCreatorID
		parameter  service.TicketAddParameter
	}
	type output struct {
		ticketID domain.TicketID
	}
	tests := []struct {
		name   string
		input  input
		output output
	}{
		{
			name: "success",
			input: input{
				operatorID: ticketCreatorID,
				parameter:  parameter,
			},
			output: output{
				ticketID: ticketID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appUserRepoMock := new(service_mock.AppUserRepository)
			ticketRepositoryMock := new(service_mock.TicketRepository)
			rfMock := new(service_mock.RepositoryFactory)
			appUserModelMock := new(domain_mock.AppUserModel)
			ticketCreator, err := service.NewTicketCreator(appUserModelMock, rfMock)
			require.NoError(t, err)
			rfMock.On("NewAppUserRepository", ctx).Return(appUserRepoMock)
			rfMock.On("NewTicketRepository", ctx).Return(ticketRepositoryMock)
			transactionManager, err := gateway.NewNoneTransactionManager(rfMock)
			require.NoError(t, err)
			usecase := usecase.NewTicketCreatorUsecase(transactionManager)

			// given
			appUserRepoMock.On("FindTicketCreatorByID", ctx, tt.input.operatorID).Return(ticketCreator, nil)
			ticketRepositoryMock.On("AddTicket", ctx, ticketCreator, tt.input.parameter).Return(tt.output.ticketID, nil)

			// when
			addedTicketID, err := usecase.AddTicket(ctx, ticketCreatorID, parameter)

			// then
			require.NoError(t, err)
			assert.Equal(t, tt.output.ticketID.Int(), addedTicketID.Int())
			appUserRepoMock.AssertCalled(t, "FindTicketCreatorByID", ctx, ticketCreatorID)
			ticketRepositoryMock.AssertCalled(t, "AddTicket", ctx, ticketCreator, parameter)
		})
	}
}
