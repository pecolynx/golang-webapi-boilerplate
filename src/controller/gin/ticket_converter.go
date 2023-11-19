package handler

import (
	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
)

func ToTicketHTTPEntity(ticket domain.TicketModel) (TicketResponseHTTPEntity, error) {
	e := TicketResponseHTTPEntity{
		BaseModel: BaseModel{
			Version:   ticket.GetVersion(),
			CreatedBy: ticket.GetCreatedBy(),
			UpdatedBy: ticket.GetUpdatedBy(),
		},
		ID:      ticket.GetTicketID().Int(),
		Name:    ticket.GetName(),
		Content: ticket.GetContent(),
	}

	if err := libdomain.Validator.Struct(e); err != nil {
		return TicketResponseHTTPEntity{}, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return e, nil
}

func ToTicketAddParameter(param *TicketAddParameter) (service.TicketAddParameter, error) {
	serviceParam, err := service.NewTicketAddParameter(param.Name, param.Description)

	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return serviceParam, nil
}
