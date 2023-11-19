package gateway

import (
	"context"
	"time"

	"gorm.io/gorm"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
)

type repositoryFactory struct {
	driverName string
	db         *gorm.DB
	location   *time.Location
}

func NewRepositoryFactory(ctx context.Context, driverName string, db *gorm.DB, location *time.Location) (service.RepositoryFactory, error) {
	if db == nil {
		return nil, liberrors.Errorf("db is nil. err: %w", libdomain.ErrInvalidArgument)
	}

	return &repositoryFactory{
		driverName: driverName,
		db:         db,
		location:   location,
	}, nil
}

func (f *repositoryFactory) NewAppUserRepository(ctx context.Context) service.AppUserRepository {
	return newAppUserRepository(ctx, f.driverName, f.db, f)
}

func (f *repositoryFactory) NewTicketRepository(ctx context.Context) service.TicketRepository {
	return newTicketRepository(ctx, f.driverName, f.db, f)
}

func (f *repositoryFactory) NewRBACRepository(ctx context.Context) service.RBACRepository {
	return NewRBACRepository(ctx, f.db)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
