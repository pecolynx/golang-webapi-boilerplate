package handler

import (
	"context"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
)

func ToTicketResponse(ticket domain.TicketModel) (*TicketResponseDetails, error) {
	e := &TicketResponseDetails{
		BaseModel: BaseModel{
			Version:   ticket.GetVersion(),
			CreatedBy: ticket.GetCreatedBy(),
			UpdatedBy: ticket.GetUpdatedBy(),
		},
		ID:          ticket.GetTicketID().Int(),
		Title:       ticket.GetTitle(),
		Description: ticket.GetDescription(),
	}

	if err := libdomain.Validator.Struct(e); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return e, nil
}

func ToTicketAddParameter(param *TicketAddParameter) (service.TicketAddParameter, error) {
	serviceParam, err := service.NewTicketAddParameter(param.Title, param.Description)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return serviceParam, nil
}

func ToTicketSearchCondition(param *TicketFindParameter) (service.TicketSearchCondition, error) {
	serviceParam, err := service.NewTicketSearchCondition(param.PageNo, param.PageSize)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return serviceParam, nil
}

func ToTicketFindResponse(ctx context.Context, result service.TicketSearchResult) (*TicketFindResponse, error) {
	tickets := make([]*TicketResponseSummary, len(result.GetResults()))
	for i, p := range result.GetResults() {
		baseModel, err := NewBaseModel(p)
		if err != nil {
			return nil, liberrors.Errorf("new Model. err: %w", err)
		}

		tickets[i] = &TicketResponseSummary{
			BaseModel:   baseModel,
			ID:          p.GetTicketID().Int(),
			Title:       p.GetTitle(),
			Description: p.GetDescription(),
		}
	}

	e := &TicketFindResponse{
		TotalCount: result.GetTotalCount(),
		Results:    tickets,
	}

	if err := libdomain.Validator.Struct(e); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return e, nil
}
