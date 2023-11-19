package domain

type TicketCreatorID interface {
	AppUserID
	Int() int
}

type ticketCreatorID struct {
	Value int
}

func NewTicketCreatorID(value int) (TicketCreatorID, error) {
	return &ticketCreatorID{
		Value: value,
	}, nil
}

func (v *ticketCreatorID) Int() int {
	return v.Value
}
