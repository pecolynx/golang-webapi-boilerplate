package service

import (
	"context"
	"errors"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
)

var ErrTicketNotFound = errors.New("ticket not found")
var ErrTicketAlreadyExists = errors.New("ticket already exists")
var ErrTicketPermissionDenied = errors.New("permission denied")

type TicketRepository interface {
	AddTicket(ctx context.Context, operatorID domain.AppUserID, param TicketAddParameter) (domain.TicketID, error)

	RemoveTicket(ctx context.Context, operatorID domain.AppUserID, ticketID domain.TicketID, version int) error

	FindMyTickets(ctx context.Context, operatorID domain.AppUserID, param TicketSearchCondition) (TicketSearchResult, error)

	CanDo(ctx context.Context, operatorID domain.AppUserID, ticketID domain.TicketID, action domain.RBACAction) (bool, error)
}

type TicketAddParameter interface {
	GetTitle() string
	GetDescription() string
}

type ticketAddParameter struct {
	Title       string
	Description string
}

func NewTicketAddParameter(title string, description string) (TicketAddParameter, error) {
	m := &ticketAddParameter{
		Title:       title,
		Description: description,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *ticketAddParameter) GetTitle() string {
	return p.Title
}

func (p *ticketAddParameter) GetDescription() string {
	return p.Description
}

type TicketSearchCondition interface {
	GetPageNo() int
	GetPageSize() int
}

type ticketSearchCondition struct {
	PageNo   int `validate:"required,gte=1"`
	PageSize int `validate:"required,gte=1,lte=1000"`
}

func NewProblemSearchCondition(pageNo, pageSize int) (TicketSearchCondition, error) {
	m := &ticketSearchCondition{
		PageNo:   pageNo,
		PageSize: pageSize,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (c *ticketSearchCondition) GetPageNo() int {
	return c.PageNo
}

func (c *ticketSearchCondition) GetPageSize() int {
	return c.PageSize
}

type TicketSearchResult interface {
	GetTotalCount() int
	GetResults() []domain.TicketModel
}

type ticketSearchResult struct {
	TotalCount int
	Results    []domain.TicketModel
}

func NewProblemSearchResult(totalCount int, results []domain.TicketModel) (TicketSearchResult, error) {
	m := &ticketSearchResult{
		TotalCount: totalCount,
		Results:    results,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *ticketSearchResult) GetTotalCount() int {
	return m.TotalCount
}

func (m *ticketSearchResult) GetResults() []domain.TicketModel {
	return m.Results
}
