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
	GetName() string
	GetContent() string
}

type ticketModel struct {
	libdomain.BaseModel
	TicketID TicketID
	Name     string `validate:"required"`
	Content  string `validate:"required"`
}

func NewTicketModel(model libdomain.BaseModel, ticketID TicketID, name string, content string) (TicketModel, error) {
	m := &ticketModel{
		BaseModel: model,
		TicketID:  ticketID,
		Name:      name,
		Content:   content,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *ticketModel) GetTicketID() TicketID {
	return m.TicketID
}

func (m *ticketModel) GetName() string {
	return m.Name
}

func (m *ticketModel) GetContent() string {
	return m.Content
}
