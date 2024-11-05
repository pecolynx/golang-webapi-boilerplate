package gateway_test

import (
	"context"
	"testing"

	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type testService struct {
	driverName string
	db         *gorm.DB
	rf         service.RepositoryFactory
}

func testNewAppUser(t *testing.T, ctx context.Context, ts testService, loginID, username, password string) service.AppUser {
	appUserRepo := ts.rf.NewAppUserRepository(ctx)
	userID1, err := appUserRepo.AddAppUser(ctx, testNewAppUserAddParameter(t, loginID, username, password))
	require.NoError(t, err)
	user1, err := appUserRepo.FindAppUserByID(ctx, owner, userID1)
	require.NoError(t, err)
	require.Equal(t, loginID, user1.GetLoginID())

	return user1
}

func testNewAppUserAddParameter(t *testing.T, loginID, username, password string) service.AppUserAddParameter {
	p, err := service.NewAppUserAddParameter(loginID, username, password)
	require.NoError(t, err)
	return p
}
