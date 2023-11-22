package domain

type StandardUserID interface {
	AppUserID
	Int() int
}

type standardUserID struct {
	Value int
}

func NewStandardUserID(value int) (StandardUserID, error) {
	return &standardUserID{
		Value: value,
	}, nil
}

func (v *standardUserID) Int() int {
	return v.Value
}
