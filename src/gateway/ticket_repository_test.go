package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/gateway"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
	testlibgateway "github.com/pecolynx/golang-webapi-boilerplate/testlib/gateway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ticketRepository_AddTicket(t *testing.T) {
	ctx := context.Background()
	loc := time.UTC

	standardUserID, err := domain.NewStandardUserID(1)
	require.NoError(t, err)
	parameter, err := service.NewTicketAddParameter("TITLE", "DESCRIPTION")
	require.NoError(t, err)
	type input struct {
	}
	type output struct {
	}
	tests := []struct {
		name   string
		input  input
		output output
	}{
		{
			name:   "",
			input:  input{},
			output: output{},
		},
	}

	for driverName, db := range testlibgateway.ListDB() {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rf, err := gateway.NewRepositoryFactory(ctx, driverName, db, loc)
				require.NoError(t, err)

				ticketRepo := rf.NewTicketRepository(ctx)
				// when
				ticketID, err := ticketRepo.AddTicket(ctx, standardUserID, parameter)
				require.NoError(t, err)

				// then
				assert.Greater(t, ticketID.Int(), 0)
				can, err := ticketRepo.CanDo(ctx, standardUserID, ticketID, domain.RBACRemoveAction)
				require.NoError(t, err)
				assert.True(t, can)

				// clean up
				err = ticketRepo.RemoveTicket(ctx, standardUserID, ticketID, 1)
				require.NoError(t, err)
			})
		}
	}
}

func Test_ticketRepository_FindMyTickets(t *testing.T) {
	ctx := context.Background()
	loc := time.UTC
	standardUserID1, err := domain.NewStandardUserID(1)
	require.NoError(t, err)
	standardUserID2, err := domain.NewStandardUserID(2)
	require.NoError(t, err)
	parameter, err := service.NewTicketAddParameter("TITLE", "DESCRIPTION")
	require.NoError(t, err)

	type input struct {
	}
	type output struct {
	}
	tests := []struct {
		name   string
		input  input
		output output
	}{
		{},
	}
	for driverName, db := range testlibgateway.ListDB() {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rf, err := gateway.NewRepositoryFactory(ctx, driverName, db, loc)
				require.NoError(t, err)

				ticketRepo := rf.NewTicketRepository(ctx)

				// user1 adds three tickets
				userTickets1 := make([]domain.TicketID, 0)
				for i := 0; i < 3; i++ {
					ticketID, err := ticketRepo.AddTicket(ctx, standardUserID1, parameter)
					require.NoError(t, err)
					assert.Greater(t, ticketID.Int(), 0)
					userTickets1 = append(userTickets1, ticketID)
				}

				// user2 adds three tickets
				userTickets2 := make([]domain.TicketID, 0)
				for i := 0; i < 3; i++ {
					ticketID, err := ticketRepo.AddTicket(ctx, standardUserID2, parameter)
					require.NoError(t, err)
					assert.Greater(t, ticketID.Int(), 0)
					userTickets2 = append(userTickets2, ticketID)
				}

				searchCondition, err := service.NewTicketSearchCondition(1, 50)
				require.NoError(t, err)
				results, err := ticketRepo.FindMyTickets(ctx, standardUserID1, searchCondition)
				require.NoError(t, err)
				tickets := results.GetResults()
				assert.Len(t, tickets, 1)

				// clean up
				for _, ticketID := range userTickets1 {
					err = ticketRepo.RemoveTicket(ctx, standardUserID1, ticketID, 1)
				}
				for _, ticketID := range userTickets2 {
					err = ticketRepo.RemoveTicket(ctx, standardUserID2, ticketID, 1)
				}
				require.NoError(t, err)

			})
		}
	}
}
