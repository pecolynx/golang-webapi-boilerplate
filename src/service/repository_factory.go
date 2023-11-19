package service

import (
	"context"
)

type RepositoryFactory interface {
	NewRBACRepository(ctx context.Context) RBACRepository
	NewTicketRepository(ctx context.Context) TicketRepository
	NewAppUserRepository(ctx context.Context) AppUserRepository
}
