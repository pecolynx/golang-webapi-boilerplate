package domain

import (
	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
)

type TicketID interface {
	Int() int
}

type ticketID struct {
	value int
}

func NewTicketID(value int) (TicketID, error) {
	return &ticketID{
		value: value,
	}, nil
}

func (v *ticketID) Int() int {
	return v.value
}

type TicketModel interface {
	libdomain.BaseModel
	GetTicketID() TicketID
	GetTitle() string
	GetDescription() string
}

type ticketModel struct {
	libdomain.BaseModel
	TicketID    TicketID
	Title       string `validate:"required"`
	Description string
}

func NewTicketModel(model libdomain.BaseModel, ticketID TicketID, title string, description string) (TicketModel, error) {
	m := &ticketModel{
		BaseModel:   model,
		TicketID:    ticketID,
		Title:       title,
		Description: description,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *ticketModel) GetTicketID() TicketID {
	return m.TicketID
}

func (m *ticketModel) GetTitle() string {
	return m.Title
}

func (m *ticketModel) GetDescription() string {
	return m.Description
}
