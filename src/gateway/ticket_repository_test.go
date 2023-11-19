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

	ticketCreatorID, err := domain.NewTicketCreatorID(1)
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
				ticketID, err := ticketRepo.AddTicket(ctx, ticketCreatorID, parameter)
				require.NoError(t, err)
				assert.Greater(t, ticketID.Int(), 0)
				can, err := ticketRepo.CanDo(ctx, ticketCreatorID, ticketID, domain.RBACRemoveAction)
				require.NoError(t, err)
				assert.True(t, can)
			})
		}
	}
}
