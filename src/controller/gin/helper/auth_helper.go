package helper

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/log"
)

func HandleSecuredFunction(c *gin.Context, fn func(ctx context.Context, logger *slog.Logger, operatorID domain.AppUserID) error, errorHandle func(ctx context.Context, logger *slog.Logger, c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	authLogger := liblog.GetLoggerFromContext(ctx, log.AppAuthLoggerContextKey)

	appUserID := c.GetInt("AuthorizedUser")
	if appUserID == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	operatorID, err := domain.NewAppUserID(appUserID)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	authLogger.InfoContext(ctx, "", slog.Int("operator_id", operatorID.Int()))

	controllerLogger := liblog.GetLoggerFromContext(ctx, log.AppControllerLoggerContextKey)
	if err := fn(ctx, controllerLogger, operatorID); err != nil {
		if handled := errorHandle(ctx, controllerLogger, c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}
