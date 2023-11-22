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

type TicketResponseDetails struct {
	BaseModel
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type TicketAddParameter struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type TicketUpdateParameter struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type TicketFindParameter struct {
	PageNo   int `json:"pageNo" binding:"required,gte=1"`
	PageSize int `json:"pageSize" binding:"required,gte=1,lte=100"`
}

type TicketResponseSummary struct {
	BaseModel
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type TicketFindResponse struct {
	TotalCount int                      `json:"totalCount" validate:"gte=0"`
	Results    []*TicketResponseSummary `json:"results" validate:"dive"`
}
type IDResponse struct {
	ID int `json:"id"`
}
