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

	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	handler "github.com/pecolynx/golang-webapi-boilerplate/src/controller/gin"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
	"github.com/pecolynx/golang-webapi-boilerplate/src/usecase"
	usecase_mock "github.com/pecolynx/golang-webapi-boilerplate/src/usecase/mocks"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })

func initTicketRouter(t *testing.T, ticketCreatorUsecase usecase.TicketCreatorUsecase, authMiddleware gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(authMiddleware)
	g := router.Group("v1")
	fn := handler.NewInitTicketRouterFunc(ticketCreatorUsecase)
	err := fn(g)
	require.NoError(t, err)
	return router
}

func newAuthMiddleware(appUserID domain.AppUserID, authenticated bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authenticated {
			c.Set("AuthorizedUser", appUserID.Int())
		}
	}
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
			ticketID, err := domain.NewTicketID(tt.input.returnTicketID)
			require.NoError(t, err)

			// given
			appUserID, err := domain.NewAppUserID(1)
			require.NoError(t, err)
			authMiddleware := newAuthMiddleware(appUserID, tt.input.authenticated)
			ticketCreatorUsecase := new(usecase_mock.TicketCreatorUsecase)
			ticketCreatorUsecase.On("AddTicket",
				anythingOfContext,
				mock.MatchedBy(func(id domain.StandardUserID) bool { return id.Int() == 1 }),
				mock.MatchedBy(func(param service.TicketAddParameter) bool {
					return param.GetTitle() == "TITLE" && param.GetDescription() == "DESCRIPTION"
				}),
			).Return(ticketID, tt.input.returnErr)
			r := initTicketRouter(t, ticketCreatorUsecase, authMiddleware)

			// when
			body, err := json.Marshal(gin.H{
				"title":       tt.input.title,
				"description": tt.input.description,
			})
			require.NoError(t, err)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/v1/ticket", bytes.NewBuffer(body))
			require.NoError(t, err)
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
