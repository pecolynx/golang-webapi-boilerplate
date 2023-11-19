package log

import (
	"context"
	"log/slog"

	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
)

const (
	LibGatewayLoggerContextKey    libdomain.ContextKey = "lib_gateway"
	AppGORMLoggerContextKey       libdomain.ContextKey = "app_gorm"
	AppServiceLoggerContextKey    libdomain.ContextKey = "app_service"
	AppGatewayLoggerContextKey    libdomain.ContextKey = "app_gateway"
	AppControllerLoggerContextKey libdomain.ContextKey = "app_controller"
	AppGinLoggerContextKey        libdomain.ContextKey = "app_gin"
	AppTraceLoggerContextKey      libdomain.ContextKey = "app_trace"
	AppAuthLoggerContextKey       libdomain.ContextKey = "app_auth"
)

var (
	LoggerKeys = []libdomain.ContextKey{
		LibGatewayLoggerContextKey,
		AppGORMLoggerContextKey,
		AppServiceLoggerContextKey,
		AppGatewayLoggerContextKey,
		AppControllerLoggerContextKey,
		AppGinLoggerContextKey,
		AppTraceLoggerContextKey,
		AppAuthLoggerContextKey,
	}
)

func InitLogger(ctx context.Context) context.Context {
	for _, key := range LoggerKeys {
		if _, ok := liblog.Loggers[key]; !ok {
			liblog.Loggers[key] = slog.New(liblog.LogHandlers[liblog.DefaultLogLevel])
		}
		ctx = context.WithValue(ctx, key, liblog.Loggers[key])
	}
	return ctx
}
