package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	// "github.com/kujilabo/cocotola/lib/log"

	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/log"
)

func NewTraceLogMiddleware(appName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sc := trace.SpanFromContext(c.Request.Context()).SpanContext()
		if !sc.TraceID().IsValid() || !sc.SpanID().IsValid() {
			return
		}
		otTraceID := sc.TraceID().String()

		ctx := c.Request.Context()
		savedCtx := ctx
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()

		ctx = liblog.WithLoggerName(ctx, log.AppTraceLoggerContextKey)
		logger := liblog.GetLoggerFromContext(ctx, log.AppTraceLoggerContextKey)
		logger.InfoContext(ctx, fmt.Sprintf("uri: %s, method: %s", c.Request.RequestURI, c.Request.Method), slog.String("request_id", otTraceID))
		// logger := log.FromContext(ctx)
		// logger.Infof("uri: %s, method: %s", c.Request.RequestURI, c.Request.Method)

		ctx, span := tracer.Start(ctx, "TraceLog")
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()
	}
}
