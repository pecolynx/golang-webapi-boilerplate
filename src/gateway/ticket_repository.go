package gateway

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"gorm.io/gorm"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	libgateway "github.com/pecolynx/golang-webapi-boilerplate/lib/gateway"
	domain "github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/gateway/casbinquery"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
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

func (e *ticketEntity) toTicketModel() (domain.TicketModel, error) {
	baseModel, err := libdomain.NewBaseModel(e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, liberrors.Errorf("libdomain.NewModel. err: %w", err)
	}

	ticketID, err := domain.NewTicketID(e.ID)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewTicketID. value: %s, err: %w", e.ID, err)
	}

	ticket, err := domain.NewTicketModel(baseModel, ticketID, e.Title, e.Description)
	if err != nil {
		return nil, liberrors.Errorf("failed to NewTicket. entity: %+v, err: %w", e, err)
	}
	return ticket, nil
}

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

func NewRBACTicketObject(ticketID domain.TicketID) domain.RBACObject {
	return domain.RBACObject(fmt.Sprintf("ticket_%d", ticketID.Int()))
}

func NewRBACTicketAssigneeRole(ticketID domain.TicketID) domain.RBACRole {
	return domain.RBACRole(fmt.Sprintf("ticket_%d_assignee", ticketID.Int()))
}

func NewRBACAppUser(appUserID domain.AppUserID) domain.RBACUser {
	return domain.RBACUser(fmt.Sprintf("user_%d", appUserID.Int()))
}

func (r *ticketRepository) AddTicket(ctx context.Context, operatorID domain.AppUserID, param service.TicketAddParameter) (domain.TicketID, error) {
	_, span := tracer.Start(ctx, "ticketRepository.AddTicket")
	defer span.End()

	ticket := ticketEntity{
		Version:     1,
		CreatedBy:   operatorID.Int(),
		UpdatedBy:   operatorID.Int(),
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
	userObject := NewRBACAppUser(operatorID)
	ticketObject := NewRBACTicketObject(ticketID)
	ticketAssignee := NewRBACTicketAssigneeRole(ticketID)

	// the ticketAssignee role can read, update, remove
	if err := rbacRepo.AddNamedPolicy(ticketAssignee, ticketObject, domain.RBACReadAction); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: read, err: %w", err)
	}
	if err := rbacRepo.AddNamedPolicy(ticketAssignee, ticketObject, domain.RBACUpdatection); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: update, err: %w", err)
	}
	if err := rbacRepo.AddNamedPolicy(ticketAssignee, ticketObject, domain.RBACRemoveAction); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedPolicy. priv: remove, err: %w", err)
	}

	// user is assigned the ticketAssignee role
	if err := rbacRepo.AddNamedGroupingPolicy(userObject, ticketAssignee); err != nil {
		return nil, liberrors.Errorf("Failed to AddNamedGroupingPolicy. err: %w", err)
	}

	return ticketID, nil
}

func (r *ticketRepository) RemoveTicket(ctx context.Context, operatorID domain.AppUserID, ticketID domain.TicketID, version int) error {
	_, span := tracer.Start(ctx, "")
	defer span.End()

	ok, err := r.CanDo(ctx, operatorID, ticketID, domain.RBACRemoveAction)
	if err != nil {
		return err
	} else if !ok {
		return service.ErrTicketPermissionDenied
	}

	if result := r.db.Where("id = ? and version = ?", ticketID.Int(), version).Delete(&ticketEntity{}); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return service.ErrTicketNotFound
		}

		return result.Error
	}

	return nil
}

func (r *ticketRepository) FindMyTickets(ctx context.Context, operatorID domain.AppUserID, param service.TicketSearchCondition) (service.TicketSearchResult, error) {
	_, span := tracer.Start(ctx, "ticketRepository.FindMyTickets")
	defer span.End()

	limit := param.GetPageSize()
	offset := (param.GetPageNo() - 1) * param.GetPageSize()
	ticketEntities := []ticketEntity{}

	objectColumnName := "name"
	subQuery, err := casbinquery.QueryObject(r.db, r.driverName, "ticket_", objectColumnName, "user_"+strconv.Itoa(operatorID.Int()), "read")
	if err != nil {
		return nil, liberrors.Errorf("casbinquery.QueryObject. err: %w", err)
	}

	if result := r.db.Model(&ticketEntity{}).
		Joins("inner join (?) AS t3 ON `ticket`.`id`= t3."+objectColumnName, subQuery).
		Order("`ticket`.`name`").Limit(limit).Offset(offset).
		Scan(&ticketEntities); result.Error != nil {
		return nil, result.Error
	}

	tickets := make([]domain.TicketModel, len(ticketEntities))
	// priv := userD.NewPrivileges([]userD.RBACAction{domain.PrivilegeRead})
	for i, e := range ticketEntities {
		t, err := e.toTicketModel()
		if err != nil {
			return nil, liberrors.Errorf("toWorkbookModel. err: %w", err)
		}
		tickets[i] = t
	}

	var count int64
	rows, err := r.db.Raw("select count(*) from workbook inner join (?) AS t3 ON `workbook`.`id`= t3."+objectColumnName, subQuery).Rows()
	if err != nil {
		return nil, liberrors.Errorf("r.db.Raw. err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c int64
		if err := rows.Scan(&c); err != nil {
			return nil, liberrors.Errorf("rows.Scan. err: %w", err)
		}
		count += c
	}

	if count > math.MaxInt32 {
		return nil, errors.New("overflow")
	}

	result, err := service.NewProblemSearchResult(int(count), tickets)
	if err != nil {
		return nil, liberrors.Errorf(". err: %w", err)
	}

	return result, nil
}

func (r *ticketRepository) getAllRolesForTicket(ticketID domain.TicketID) []domain.RBACRole {
	return []domain.RBACRole{
		NewRBACTicketAssigneeRole(ticketID),
	}
}

func (r *ticketRepository) CanDo(ctx context.Context, operatorID domain.AppUserID, ticketID domain.TicketID, action domain.RBACAction) (bool, error) {
	rbacRepo := r.rf.NewRBACRepository(ctx)

	roleObjects := r.getAllRolesForTicket(ticketID)
	userObject := NewRBACAppUser(operatorID)
	e, err := rbacRepo.NewEnforcerWithRolesAndUsers(roleObjects, []domain.RBACUser{userObject})
	if err != nil {
		return false, liberrors.Errorf("failed to NewEnforcerWithRolesAndUsers. err: %w", err)
	}

	ticketObject := NewRBACTicketObject(ticketID)

	ok, err := e.Enforce(string(userObject), string(ticketObject), string(action))
	if err != nil {
		return false, liberrors.Errorf("e.Enforce. err: %w", err)
	}

	return ok, nil
}
