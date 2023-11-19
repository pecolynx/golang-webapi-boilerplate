package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/swag/example/basic/docs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	libconfig "github.com/pecolynx/golang-webapi-boilerplate/lib/config"
	libdomain "github.com/pecolynx/golang-webapi-boilerplate/lib/domain"
	liberrors "github.com/pecolynx/golang-webapi-boilerplate/lib/errors"
	libgateway "github.com/pecolynx/golang-webapi-boilerplate/lib/gateway"
	liblog "github.com/pecolynx/golang-webapi-boilerplate/lib/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/config"
	controller "github.com/pecolynx/golang-webapi-boilerplate/src/controller/gin"
	"github.com/pecolynx/golang-webapi-boilerplate/src/gateway"
	"github.com/pecolynx/golang-webapi-boilerplate/src/log"
	"github.com/pecolynx/golang-webapi-boilerplate/src/service"
	"github.com/pecolynx/golang-webapi-boilerplate/src/sqls"
)

const readHeaderTimeout = time.Duration(30) * time.Second

var (
	y = 2
)

func getValue(values ...string) string {
	for _, v := range values {
		if len(v) != 0 {
			return v
		}
	}
	return ""
}

func main() {
	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	appEnv := getValue(*env, os.Getenv("APP_ENV"), "local")
	slog.InfoContext(ctx, fmt.Sprintf("env: %s", appEnv))

	liberrors.UseXerrorsErrorf()

	cfg, db, sqlDB, tp := initialize(ctx, appEnv)
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	ctx = log.InitLogger(ctx)
	logger := liblog.GetLoggerFromContext(ctx, libdomain.ContextKey(cfg.App.Name))

	rff := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(ctx, cfg.DB.DriverName, db, time.UTC) // nolint:wrapcheck
	}

	appTransactionManager := initTransactionManager(db, rff)

	// baseModel, err := usermodel.NewBaseModel(1, time.Now(), time.Now(), 0, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// admin, err := usermodel.NewAppUserModel(baseModel, usermodel.AppUserID(1), "admin", "Admin")
	// if err != nil {
	// 	panic(err)
	// }
	// documentWriter, err := model.NewDocumentWriterModel(admin)
	// if err != nil {
	// 	panic(err)
	// }

	// var documentModelList []model.DocumentModel
	// appTransactionManager.Do(ctx, func(rf service.RepositoryFactory) error {
	// 	documentRepo := rf.NewDocumentRepository(ctx)
	// 	searchCondition, err := model.NewDocumentSearchCondition(1, 100)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	documentSearchResult, err := documentRepo.FindPrivateDocuments(ctx, documentWriter, searchCondition)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	documentModelList = documentSearchResult.GetResults()
	// 	return nil
	// })

	// fmt.Println(documentModelList)

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second

	result := run(context.Background(), cfg, appTransactionManager)

	time.Sleep(gracefulShutdownTime2)
	logger.InfoContext(ctx, "exited")
	os.Exit(result)
}

func a(x int) int {
	return x * y
}

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *sdktrace.TracerProvider) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		panic(err)
	}

	// init log
	if err := libconfig.InitLog(cfg.Log); err != nil {
		panic(err)
	}

	// init tracer
	tp, err := libconfig.InitTracerProvider(cfg.App.Name, cfg.Trace)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// init db
	db, sqlDB, err := libconfig.InitDB(cfg.DB, sqls.SQL)
	if err != nil {
		panic(err)
	}

	return cfg, db, sqlDB, tp
}
func initTransactionManager(db *gorm.DB, rff gateway.RepositoryFactoryFunc) service.TransactionManager {
	appTransactionManager, err := gateway.NewTransactionManager(db, rff)
	if err != nil {
		panic(err)
	}

	return appTransactionManager
}

func run(ctx context.Context, cfg *config.Config, transactionManager service.TransactionManager) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	eg.Go(func() error {
		return appServer(ctx, cfg) // nolint:wrapcheck
	})
	eg.Go(func() error {
		return libgateway.MetricsServerProcess(ctx, cfg.App.MetricsPort, cfg.Shutdown.TimeSec1) // nolint:wrapcheck
	})
	eg.Go(func() error {
		return libgateway.SignalWatchProcess(ctx) // nolint:wrapcheck
	})
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err() // nolint:wrapcheck
	})

	if err := eg.Wait(); err != nil {
		logrus.Error(err)
		return 1
	}
	return 0
}

func appServer(ctx context.Context, cfg *config.Config) error {
	logger := liblog.GetLoggerFromContext(ctx, libdomain.ContextKey(cfg.App.Name))
	// // cors
	corsConfig := libconfig.InitCORS(cfg.CORS)
	// logrus.Infof("cors: %+v", corsConfig)

	// if err := corsConfig.Validate(); err != nil {
	// 	return liberrors.Errorf("corsConfig.Validate. err: %w", err)
	// }

	// studyMonitor := service.NewStudyMonitor()
	// studyStatUpdater := studyStatUpdater{
	// 	systemOwnerModel: systemOwnerModel,
	// 	appTransaction:   appTransaction,
	// }
	// if err := studyMonitor.Attach(&studyStatUpdater); err != nil {
	// 	return liberrors.Errorf(". err: %w", err)
	// }

	// privateRouterGroupFunc := []controller.InitRouterGroupFunc{
	// 	controller.NewInitWorkbookRouterFunc(studentUsecaseWorkbook),
	// 	controller.NewInitProblemRouterFunc(studentUsecaseProblem, newIteratorFunc),
	// 	controller.NewInitStudyRouterFunc(studentUseCaseStudy),
	// 	controller.NewInitAudioRouterFunc(studentUsecaseAudio),
	// 	controller.NewInitStatRouterFunc(studentUsecaseStat),
	// }

	publicRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitTestRouterFunc(),
	}
	router, err := controller.NewAppRouter(ctx,
		publicRouterGroupFunc,
		//privateRouterGroupFunc, pluginRouterGroupFunc, authTokenManager,
		corsConfig, cfg.App,
		//cfg.Auth,
		cfg.Debug)
	if err != nil {
		panic(err)
	}

	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docs.SwaggerInfo.Title = cfg.App.Name
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = cfg.Swagger.Host
		docs.SwaggerInfo.Schemes = []string{cfg.Swagger.Schema}
	}

	httpServer := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.App.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logger.InfoContext(ctx, fmt.Sprintf("http server listening at %v", httpServer.Addr))

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.InfoContext(ctx, fmt.Sprintf("failed to ListenAndServe. err: %v", err))
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.InfoContext(ctx, fmt.Sprintf("Server forced to shutdown. err: %v", err))
			return liberrors.Errorf(". err: %w", err)
		}
		return nil
	case err := <-errCh:
		return liberrors.Errorf(". err: %w", err)
	}
}
