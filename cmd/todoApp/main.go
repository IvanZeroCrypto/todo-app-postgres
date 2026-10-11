package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "gitverse.ru/IvanVerse/go-todoApp/internal/core/config"
	core_logger "gitverse.ru/IvanVerse/go-todoApp/internal/core/logger"
	core_postgres_pool "gitverse.ru/IvanVerse/go-todoApp/internal/core/repository/postgres/pool"
	core_http_middleware "gitverse.ru/IvanVerse/go-todoApp/internal/core/transport/http/middleware"
	core_http_server "gitverse.ru/IvanVerse/go-todoApp/internal/core/transport/http/server"
	statistics_postgres_repository "gitverse.ru/IvanVerse/go-todoApp/internal/features/statistics/repository/postgres"
	statistics_service "gitverse.ru/IvanVerse/go-todoApp/internal/features/statistics/service"
	statistics_transport_http "gitverse.ru/IvanVerse/go-todoApp/internal/features/statistics/transport/http"
	tasks_postgres_repository "gitverse.ru/IvanVerse/go-todoApp/internal/features/tasks/repository/postgres"
	tasks_service "gitverse.ru/IvanVerse/go-todoApp/internal/features/tasks/service"
	tasks_transport_http "gitverse.ru/IvanVerse/go-todoApp/internal/features/tasks/transport/http"
	users_postgres_repository "gitverse.ru/IvanVerse/go-todoApp/internal/features/users/repository/postgres"
	users_service "gitverse.ru/IvanVerse/go-todoApp/internal/features/users/service"
	users_transport_http "gitverse.ru/IvanVerse/go-todoApp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())

	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}

	defer logger.Close()
	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initialized postgres connection poll")
	pool, err := core_postgres_pool.NewConnectionPool(core_postgres_pool.NewConfigMust(), ctx)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initiaize feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initiaize feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initiaize feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPhandler(statisticsService)

	logger.Debug("initialized HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouter(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouter(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouter(statisticsTransportHTTP.Routes()...)

	// apiVersionRouterv2 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion2,
	// 	core_http_middleware.Dummy("api va middleware"))
	// apiVersionRouterv2.RegisterRouter(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter) // apiVersionRouterv2,

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}

}
