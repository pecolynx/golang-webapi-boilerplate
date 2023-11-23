package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/config"
	"github.com/pecolynx/golang-webapi-boilerplate/src/controller/auth"
	handler "github.com/pecolynx/golang-webapi-boilerplate/src/controller/gin"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
	"github.com/pecolynx/golang-webapi-boilerplate/src/usecase"
	usecase_mock "github.com/pecolynx/golang-webapi-boilerplate/src/usecase/mocks"
)

var (
	anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })
	corsConfig        cors.Config
	appConfig         *config.AppConfig
	authConfig        *config.AuthConfig
	debugConfig       *config.DebugConfig
	authTokenManager  auth.AuthTokenManager
)

func init() {
	corsConfig = cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}
	appConfig = &config.AppConfig{
		Name:        "test",
		HTTPPort:    8080,
		MetricsPort: 8081,
	}
	authConfig = &config.AuthConfig{
		SigningKey:          "ah5T9Y9V2JPU74fhCtHQfDqLp3Zg8ZNc",
		AccessTokenTTLMin:   1,
		RefreshTokenTTLHour: 1,
	}
	debugConfig = &config.DebugConfig{
		Gin:  false,
		Wait: false,
	}

	signingKey := []byte(authConfig.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager = auth.NewAuthTokenManager(signingKey, signingMethod, time.Duration(authConfig.AccessTokenTTLMin)*time.Minute, time.Duration(authConfig.RefreshTokenTTLHour)*time.Hour)
}

func initTicketRouter(t *testing.T, ctx context.Context, ticketCreatorUsecase usecase.TicketCreatorUsecase) *gin.Engine {
	// router := gin.New()
	// router.Use(authMiddleware)
	// g := router.Group("v1")
	fn := handler.NewInitTicketRouterFunc(ticketCreatorUsecase)
	// err := fn(g)
	// require.NoError(t, err)
	router, err := handler.NewAppRouter(ctx, []handler.InitRouterGroupFunc{fn}, nil, corsConfig, appConfig, authConfig, debugConfig)
	require.NoError(t, err)
	return router
}

func parseJSON(t *testing.T, b *bytes.Buffer) interface{} {
	respBytes, err := io.ReadAll(b)
	require.NoError(t, err)
	obj, err := oj.Parse(respBytes)
	require.NoError(t, err)
	return obj
}

func parseExpr(t *testing.T, v string) jp.Expr {
	expr, err := jp.ParseString(v)
	require.NoError(t, err)
	return expr
}

func createTokenSet(t *testing.T, ctx context.Context) *auth.TokenSet {
	appUserID, err := domain.NewAppUserID(1)
	require.NoError(t, err)

	now := time.Now()
	baseModel, err := libdomain.NewBaseModel(1, now, now, 1, 1)
	require.NoError(t, err)

	appUserModel, err := domain.NewAppUserModel(baseModel, appUserID, "LOGIN_ID", "USERNAME")
	require.NoError(t, err)

	tokenSet, err := authTokenManager.CreateTokenSet(ctx, appUserModel)
	require.NoError(t, err)

	return tokenSet
}

func Test_ticketHandler_AddTicket(t *testing.T) {
	ctx := context.Background()
	type input struct {
		authenticated  bool
		title          string
		description    string
		returnTicketID int
		returnErr      error
	}
	type output struct {
		statusCode int
		message    string
	}
	tests := []struct {
		name   string
		input  input
		output output
	}{
		{
			name: "ok",
			input: input{
				authenticated:  true,
				title:          "TITLE",
				description:    "DESCRIPTION",
				returnTicketID: 456,
			},
			output: output{
				statusCode: 200,
			},
		},
		{
			name: "unauthorized",
			input: input{
				authenticated: false,
			},
			output: output{statusCode: 401},
		},
		{
			name: "bad_request",
			input: input{
				authenticated: true,
				title:         "",
				description:   "DESCRIPTION",
			},
			output: output{
				statusCode: 400,
				message:    "Key: 'TicketAddParameter.Title' Error:Field validation for 'Title' failed on the 'required' tag",
			},
		},
		{
			name: "internal_server_error",
			input: input{
				authenticated: true,
				title:         "TITLE",
				description:   "DESCRIPTION",
				returnErr:     errors.New("ERROR"),
			},
			output: output{
				statusCode: 500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ticketID, err := domain.NewTicketID(tt.input.returnTicketID)
			require.NoError(t, err)
			ticketCreatorUsecase := new(usecase_mock.TicketCreatorUsecase)
			ticketCreatorUsecase.On("AddTicket",
				anythingOfContext,
				mock.MatchedBy(func(id domain.StandardUserID) bool { return id.Int() == 1 }),
				mock.MatchedBy(func(param service.TicketAddParameter) bool {
					return param.GetTitle() == "TITLE" && param.GetDescription() == "DESCRIPTION"
				}),
			).Return(ticketID, tt.input.returnErr)
			r := initTicketRouter(t, ctx, ticketCreatorUsecase)

			// when
			body, err := json.Marshal(gin.H{
				"title":       tt.input.title,
				"description": tt.input.description,
			})
			require.NoError(t, err)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/v1/ticket", bytes.NewBuffer(body))
			require.NoError(t, err)
			if tt.input.authenticated {
				tokenSet := createTokenSet(t, ctx)
				req.Header.Set("Authorization", "Bearer "+tokenSet.AccessToken)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// then
			jsonObj := parseJSON(t, w.Body)
			assert.Equal(t, tt.output.statusCode, w.Code)

			// - check the message
			if len(tt.output.message) != 0 {
				messageExpr := parseExpr(t, "$.message")
				message := messageExpr.Get(jsonObj)
				assert.Equal(t, tt.output.message, message[0])
			}

			// - check the status code
			if tt.input.returnTicketID != 0 {
				idExpr := parseExpr(t, "$.id")
				id := idExpr.Get(jsonObj)
				assert.Equal(t, int64(tt.input.returnTicketID), id[0])
			}

			if 400 <= w.Code && w.Code < 500 {
				return
			}

			ticketCreatorUsecase.AssertExpectations(t)
		})
	}
}
