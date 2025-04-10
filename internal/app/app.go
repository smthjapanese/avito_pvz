package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/smthjapanese/avito_pvz/internal/config"
	"github.com/smthjapanese/avito_pvz/internal/pkg/database"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
)

type App struct {
	cfg           *config.Config
	httpServer    *http.Server
	grpcServer    *grpc.Server
	metricsServer *http.Server
	logger        logger.Logger
	metrics       *metrics.Metrics
	db            *database.Database
}

func NewApp(cfg *config.Config) (*App, error) {
	// Инициализация логгера
	l, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Инициализация метрик
	m := metrics.NewMetrics()

	// Подключение к базе данных
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Создание HTTP сервера
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()

	// Создание сервера для метрик
	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	metricsServer := &http.Server{
		Addr:    ":" + cfg.Server.MetricsPort,
		Handler: metricsRouter,
	}

	return &App{
		cfg:           cfg,
		httpServer:    httpServer,
		grpcServer:    grpcServer,
		metricsServer: metricsServer,
		logger:        l,
		metrics:       m,
		db:            db,
	}, nil
}

// Run запускает приложение
func (a *App) Run() error {
	// TODO: Инициализация репозиториев, use case, обработчиков

	go func() {
		a.logger.Info(fmt.Sprintf("Starting HTTP server on port %s", a.cfg.Server.HTTPPort))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal(fmt.Sprintf("Failed to start HTTP server: %v", err))
		}
	}()

	go func() {
		a.logger.Info(fmt.Sprintf("Starting gRPC server on port %s", a.cfg.Server.GRPCPort))
		lis, err := net.Listen("tcp", ":"+a.cfg.Server.GRPCPort)
		if err != nil {
			a.logger.Fatal(fmt.Sprintf("Failed to listen for gRPC: %v", err))
		}
		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Fatal(fmt.Sprintf("Failed to start gRPC server: %v", err))
		}
	}()

	go func() {
		a.logger.Info(fmt.Sprintf("Starting metrics server on port %s", a.cfg.Server.MetricsPort))
		if err := a.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal(fmt.Sprintf("Failed to start metrics server: %v", err))
		}
	}()

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	// Остановка HTTP сервера
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := a.metricsServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown metrics server: %w", err)
	}

	a.grpcServer.GracefulStop()

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
	}

	return nil
}
