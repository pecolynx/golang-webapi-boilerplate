package handler

import (
	"time"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
)

type BaseModel struct {
	Version   int       `json:"version" validate:"required,gte=1"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy int       `json:"createdBy" validate:"gte=0"`
	UpdatedBy int       `json:"updatedBy" validate:"gte=0"`
}

func NewBaseModel(model libdomain.BaseModel) (BaseModel, error) {
	m := BaseModel{
		Version:   model.GetVersion(),
		CreatedAt: model.GetCreatedAt(),
		UpdatedAt: model.GetUpdatedAt(),
		CreatedBy: model.GetCreatedBy(),
		UpdatedBy: model.GetUpdatedBy(),
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return BaseModel{}, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

type TicketResponseHTTPEntity struct {
	BaseModel
	ID      int    `json:"id"`
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type TicketAddParameter struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type TicketUpdateParameter struct {
	Name string `json:"name" binding:"required"`
}

type IDResponse struct {
	ID int `json:"id"`
}
