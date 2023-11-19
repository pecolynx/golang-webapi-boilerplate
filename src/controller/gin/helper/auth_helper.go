package helper

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/domain"
	"github.com/pecolynx/golang-webapi-boilerplate/src/log"
)

func HandleSecuredFunction(c *gin.Context, fn func(operatorID domain.AppUserID) error, errorHandle func(c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	logger := liblog.GetLoggerFromContext(ctx, log.AppAuthLoggerContextKey)
	operatorID, err := domain.NewAppUserID(c.GetInt("AuthorizedUser"))
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	logger.InfoContext(ctx, "", slog.Int("operator_id", operatorID.Int()))
	if err := fn(operatorID); err != nil {
		if handled := errorHandle(c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}
