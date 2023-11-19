package gateway

import (
	"time"

	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"

	// 	"time"

	// 	"github.com/kujilabo/cocotola/cocotola-api/src/user/domain"
	// 	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
)

type BaseModelEntity struct {
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int
	UpdatedBy int
}

func (e *BaseModelEntity) toModel() (libdomain.BaseModel, error) {
	model, err := libdomain.NewBaseModel(e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
	if err != nil {
		return nil, liberrors.Errorf("libdomain.NewBaseModel. err: %w", err)
	}

	return model, nil
}

// type JunctionModelEntity struct {
// 	CreatedAt time.Time
// 	CreatedBy uint
// }

// // func (e *junctionModelEntity) toModel() (domain.Model, error) {
// // 	return domain.NewModel(e.ID, e.Version, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy)
// // }
