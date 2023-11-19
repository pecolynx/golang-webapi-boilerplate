package domain

import (
	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
)

type AppUserID interface {
	Int() int
}

type appUserID struct {
	Value int
}

func NewAppUserID(value int) (AppUserID, error) {
	return &appUserID{
		Value: value,
	}, nil
}

func (v *appUserID) Int() int {
	return v.Value
}

type AppUserModel interface {
	libdomain.BaseModel
	GetAppUserID() AppUserID
	GetLoginID() string
	GetUsername() string
}

type appUserModel struct {
	libdomain.BaseModel
	AppUserID AppUserID
	LoginID   string `validate:"required"`
	Username  string `validate:"required"`
}

func NewAppUserModel(model libdomain.BaseModel, appUserID AppUserID, loginID, username string) (AppUserModel, error) {
	m := &appUserModel{
		BaseModel: model,
		AppUserID: appUserID,
		LoginID:   loginID,
		Username:  username,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (m *appUserModel) GetAppUserID() AppUserID {
	return m.AppUserID
}

func (m *appUserModel) GetLoginID() string {
	return m.LoginID
}

func (m *appUserModel) GetUsername() string {
	return m.Username
}
