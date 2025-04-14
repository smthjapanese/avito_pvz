package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbv1 "github.com/smthjapanese/avito_pvz/github.com/avito_pvz/pvz/pvz_v1"
	"github.com/smthjapanese/avito_pvz/internal/config"
	"github.com/smthjapanese/avito_pvz/internal/delivery/http/handler"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/database"
	"github.com/smthjapanese/avito_pvz/internal/pkg/jwt"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
	"github.com/smthjapanese/avito_pvz/internal/repository"
	implUsecase "github.com/smthjapanese/avito_pvz/internal/usecase"
)

type App struct {
	cfg           *config.Config
	httpServer    *http.Server
	grpcServer    *grpc.Server
	metricsServer *http.Server
	logger        logger.Logger
	metrics       *metrics.Metrics
	db            *database.Database
	repositories  *repository.Repositories
	useCases      *implUsecase.UseCases
	tokenManager  *jwt.Manager
	httpHandler   *handler.Handler
}

// GetPVZUseCase возвращает PVZ use case
func (a *App) GetPVZUseCase() usecase.PVZUseCase {
	return a.useCases.PVZ
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

	// Инициализация JWT менеджера
	tokenManager := jwt.NewManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiration)

	// Инициализация репозиториев
	repos := repository.NewRepositories(db)

	// Инициализация use cases
	useCases := implUsecase.NewUseCases(repos, tokenManager)

	// Инициализация HTTP-сервера
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Инициализация обработчиков
	httpHandler := handler.NewHandler(useCases, l, m)
	httpHandler.Init(router)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	pvzServer := &PVZServer{pvzUseCase: useCases.PVZ}
	pbv1.RegisterPVZServiceServer(grpcServer, pvzServer)

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
		repositories:  repos,
		useCases:      useCases,
		tokenManager:  tokenManager,
		httpHandler:   httpHandler,
	}, nil
}

// Run запускает приложение
func (a *App) Run() error {
	// Запуск HTTP-сервера
	go func() {
		a.logger.Info(fmt.Sprintf("Starting HTTP server on port %s", a.cfg.Server.HTTPPort))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal(fmt.Sprintf("Failed to start HTTP server: %v", err))
		}
	}()

	// Запуск gRPC сервера
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

	// Запуск сервера метрик
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

// PVZServer реализует gRPC сервер для PVZ
type PVZServer struct {
	pbv1.UnimplementedPVZServiceServer
	pvzUseCase usecase.PVZUseCase
}

// GetPVZList реализует gRPC метод для получения списка ПВЗ
func (s *PVZServer) GetPVZList(ctx context.Context, req *pbv1.GetPVZListRequest) (*pbv1.GetPVZListResponse, error) {
	pvzs, err := s.pvzUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := &pbv1.GetPVZListResponse{}
	for _, pvz := range pvzs {
		response.Pvzs = append(response.Pvzs, &pbv1.PVZ{
			Id:               pvz.ID.String(),
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             string(pvz.City),
		})
	}

	return response, nil
}
