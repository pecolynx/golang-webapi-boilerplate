package gateway

import (
	"context"
	"errors"

	"gorm.io/gorm"

	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	libgateway "github.com/pecolynx/golang-webapi-boilerplate/lib/gateway"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
)

var (
	AppUserTableName = "app_user"

	// SystemOwnerLoginID   = "system-owner"
	// SystemStudentLoginID = "system-student"
	// GuestLoginID         = "guest"

	// AdministratorRole = "Administrator"
	// OwnerRole         = "Owner"
	// ManagerRole       = "Manager"
	// UserRole          = "User"
	// GuestRole         = "Guest"
	// UnknownRole       = "Unknown"
)

type appUserRepository struct {
	db *gorm.DB
	rf service.RepositoryFactory
}

type appUserEntity struct {
	BaseModelEntity
	ID             int
	LoginID        string
	Username       string
	HashedPassword string
}

func (e *appUserEntity) TableName() string {
	return AppUserTableName
}

func (e *appUserEntity) toAppUserModel() (domain.AppUserModel, error) {
	model, err := e.toModel()
	if err != nil {
		return nil, liberrors.Errorf("e.toModel. err: %w", err)
	}

	appUserID, err := domain.NewAppUserID(e.ID)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	appUserModel, err := domain.NewAppUserModel(model, appUserID, e.LoginID, e.Username)
	if err != nil {
		return nil, liberrors.Errorf("domain.NewAppUserModel. err: %w", err)
	}

	return appUserModel, nil
}

// func (e *appUserEntity) toAppUser(ctx context.Context, rf service.RepositoryFactory) (service.AppUser, error) {
// 	appUserModel, err := e.toAppUserModel()
// 	if err != nil {
// 		return nil, liberrors.Errorf("e.toAppUserModel. err: %w", err)
// 	}

// 	appUser, err := service.NewAppUser(ctx, rf, appUserModel)
// 	if err != nil {
// 		return nil, liberrors.Errorf("service.NewAppUser. err: %w", err)
// 	}

// 	return appUser, nil
// }

// func NewAppUserRepository(ctx context.Context, rf service.RepositoryFactory, db *gorm.DB) service.AppUserRepository {
// 	if rf == nil {
// 		panic(errors.New("rf is nil"))
// 	} else if db == nil {
// 		panic(errors.New("db is nil"))
// 	}
// 	return &appUserRepository{
// 		rf: rf,
// 		db: db,
// 	}
// }

func newAppUserRepository(ctx context.Context, driverName string, db *gorm.DB, rf service.RepositoryFactory) service.AppUserRepository {
	return &appUserRepository{
		db: db,
		rf: rf,
	}
}
func (r *appUserRepository) FindTicketCreatorByID(ctx context.Context, standardUserID domain.StandardUserID) (service.TicketCreator, error) {
	_, span := tracer.Start(ctx, "appUserRepository.FindAppUserByLoginID")
	defer span.End()

	appUser := appUserEntity{}
	if result := r.db.Where(&appUserEntity{
		ID: standardUserID.Int(),
	}).First(&appUser); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrAppUserNotFound
		}

		return nil, result.Error
	}

	appUserModel, err := appUser.toAppUserModel()
	if err != nil {
		return nil, err
	}
	return service.NewTicketCreator(appUserModel, r.rf)
}

func (r *appUserRepository) AddAppUser(ctx context.Context, operator domain.OwnerModel, param service.AppUserAddParameter) (domain.AppUserID, error) {
	_, span := tracer.Start(ctx, "appUserRepository.AddAppUser")
	defer span.End()

	hashedPassword, err := libgateway.HashPassword(param.GetPassword())
	if err != nil {
		return nil, liberrors.Errorf("libgateway.HashPassword. err: %w", err)
	}

	appUserEntity := appUserEntity{
		BaseModelEntity: BaseModelEntity{
			Version:   1,
			CreatedBy: operator.GetAppUserID().Int(),
			UpdatedBy: operator.GetAppUserID().Int(),
		},
		LoginID:        param.GetLoginID(),
		Username:       param.GetUsername(),
		HashedPassword: hashedPassword,
	}

	if result := r.db.Create(appUserEntity); result.Error != nil {
		return nil, liberrors.Errorf("db.Create. err: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrAppUserAlreadyExists))
	}

	appUserID, err := domain.NewAppUserID(appUserEntity.ID)
	if err != nil {
		return nil, err
	}

	return appUserID, nil
}
