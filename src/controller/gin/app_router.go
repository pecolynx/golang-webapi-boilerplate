package handler

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/config"
	"github.com/pecolynx/golang-webapi-boilerplate/src/controller/gin/middleware"
	"github.com/pecolynx/golang-webapi-boilerplate/src/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/usecase"
)

type InitRouterGroupFunc func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error

func NewInitTestRouterFunc() InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		test := parentRouterGroup.Group("test")
		for _, m := range middleware {
			test.Use(m)
		}
		test.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		return nil
	}
}

func NewInitTicketRouterFunc(ticketCreatorUsecase usecase.TicketCreatorUsecase) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		ticket := parentRouterGroup.Group("ticket")
		ticketHandler := NewTicketHandler(ticketCreatorUsecase)
		ticket.POST("", ticketHandler.AddTicket)
		return nil
	}
}

func NewAppRouter(ctx context.Context, initPublicRouterFunc []InitRouterGroupFunc, initPrivateRouterFunc []InitRouterGroupFunc, //authTokenManager service.AuthTokenManager,
	corsConfig cors.Config, appConfig *config.AppConfig,
	// authConfig *config.AuthConfig,
	debugConfig *config.DebugConfig) (*gin.Engine, error) {
	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	if debugConfig.GinMode {
		ginLogger := liblog.GetLoggerFromContext(ctx, log.AppGinLoggerContextKey)
		router.Use(sloggin.New(ginLogger))
	}

	if debugConfig.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(appConfig.Name))
		v1.Use(middleware.NewTraceLogMiddleware(appConfig.Name))

		for _, fn := range initPublicRouterFunc {
			if err := fn(v1); err != nil {
				return nil, err
			}
		}
	}

	return router, nil
}
