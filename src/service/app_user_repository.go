package service

import (
	"context"
	"errors"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
)

var ErrAppUserNotFound = errors.New("appUser not found")
var ErrAppUserAlreadyExists = errors.New("appUser already exists")
var ErrAppUserPermissionDenied = errors.New("permission denied")

type AppUserRepository interface {
	FindTicketCreatorByID(ctx context.Context, standardUserID domain.StandardUserID) (TicketCreator, error)
}

type AppUserAddParameter interface {
	GetLoginID() string
	GetUsername() string
	GetPassword() string
}

type appUserAddParameter struct {
	LoginID  string
	Username string
	Password string
}

func NewAppUserAddParameter(loginID, username, password string) (AppUserAddParameter, error) {
	m := &appUserAddParameter{
		LoginID:  loginID,
		Username: username,
		Password: password,
	}

	if err := libdomain.Validator.Struct(m); err != nil {
		return nil, liberrors.Errorf("libdomain.Validator.Struct. err: %w", err)
	}

	return m, nil
}

func (p *appUserAddParameter) GetLoginID() string {
	return p.LoginID
}

func (p *appUserAddParameter) GetUsername() string {
	return p.Username
}

func (p *appUserAddParameter) GetPassword() string {
	return p.Password
}
