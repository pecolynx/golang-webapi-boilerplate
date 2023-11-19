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
