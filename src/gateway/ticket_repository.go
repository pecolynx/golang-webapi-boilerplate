package gateway

import (
	"context"
	"fmt"
	"time"

	// casbinquery "github.com/pecolynx/casbin-query"

	// "github.com/kujilabo/cocotola/cocotola-api/src/app/gateway/casbinquery"

	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	libgateway "github.com/pecolynx/golang-webapi-boilerplate/lib/gateway"
	domain "github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
	"gorm.io/gorm"
	// "github.com/kujilabo/cocotola/cocotola-api/src/app/gateway/casbinquery"
)

type ticketEntity struct {
	ID          int
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   int
	UpdatedBy   int
	Title       string
	Description string
}

func (e *ticketEntity) TableName() string {
	return "ticket"
}

// func (e *ticketEntity) toTicketModel() (domain.TicketModel, error) {
// 	baseModel, err := libdomain.NewBaseModel(e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
// 	if err != nil {
// 		return nil, liberrors.Errorf("libdomain.NewModel. err: %w", err)
// 	}

// 	ticketID, err := domain.NewTicketID(e.ID)
// 	if err != nil {
// 		return nil, liberrors.Errorf("failed to NewTicketID. value: %s, err: %w", e.ID, err)
// 	}

// 	ticket, err := domain.NewTicketModel(baseModel, ticketID, e.Title, e.Description)
// 	if err != nil {
// 		return nil, liberrors.Errorf("failed to NewTicket. entity: %+v, err: %w", e, err)
// 	}
// 	return ticket, nil
// }

type ticketRepository struct {
	driverName string
	db         *gorm.DB
	rf         service.RepositoryFactory
}

func newTicketRepository(ctx context.Context, driverName string, db *gorm.DB, rf service.RepositoryFactory) service.TicketRepository {
	return &ticketRepository{
		driverName: driverName,
		db:         db,
		rf:         rf,
	}
}

func NewTicketObject(ticketID domain.TicketID) domain.RBACObject {
	return domain.RBACObject(fmt.Sprintf("space_%d", ticketID.Int()))
}

func NewTicketAssignee(ticketID domain.TicketID) domain.RBACRole {
	return domain.RBACRole(fmt.Sprintf("ticket_%d_assignee", ticketID.Int()))
}

func NewAppUserObject(appUserID domain.AppUserID) domain.RBACUser {
	return domain.RBACUser(fmt.Sprintf("user_%d", appUserID.Int()))
}

func (r *ticketRepository) AddTicket(ctx context.Context, operator service.TicketCreator, param service.TicketAddParameter) (domain.TicketID, error) {
	_, span := tracer.Start(ctx, "ticketRepository.AddTicket")
	defer span.End()

	ticket := ticketEntity{
		Version:     1,
		CreatedBy:   operator.GetAppUserID().Int(),
		UpdatedBy:   operator.GetAppUserID().Int(),
		Title:       param.GetTitle(),
		Description: param.GetDescription(),
	}
	if result := r.db.Create(&ticket); result.Error != nil {
		return nil, liberrors.Errorf(". err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrTicketAlreadyExists))
	}

	ticketID, err := domain.NewTicketID(ticket.ID)
	if err != nil {
		return nil, err
	}

	rbacRepo := r.rf.NewRBACRepository(ctx)
	userObject := NewAppUserObject(operator.GetAppUserID())
	ticketObject := NewTicketObject(ticketID)
	ticketAssignee := NewTicketAssignee(ticketID)

	// the ticketAssignee role can read, update, remove
	if err := rbacRepo.AddNamedPolicy(ticketAssignee, ticketObject, domain.PrivilegeRead); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}
	if err := rbacRepo.AddNamedPolicy(ticketAssignee, ticketObject, domain.PrivilegeUpdate); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: update, err: %w", err)
	}
	if err := rbacRepo.AddNamedPolicy(ticketAssignee, ticketObject, domain.PrivilegeRemove); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: remove, err: %w", err)
	}

	// user is assigned the ticketAssignee role
	if err := rbacRepo.AddNamedGroupingPolicy(userObject, ticketAssignee); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedGroupingPolicy. err: %w", err)
	}

	return ticketID, nil
}
