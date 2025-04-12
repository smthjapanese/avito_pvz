package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smthjapanese/avito_pvz/internal/delivery/http/middleware"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
)

// Handler представляет обработчик HTTP запросов
type Handler struct {
	userHandler      *UserHandler
	pvzHandler       *PVZHandler
	receptionHandler *ReceptionHandler
	productHandler   *ProductHandler
	authMiddleware   *middleware.AuthMiddleware
	logger           logger.Logger
	metrics          metrics.MetricsInterface
}

// NewHandler создает новый экземпляр Handler
func NewHandler(
	useCases *usecase.UseCases,
	logger logger.Logger,
	metrics metrics.MetricsInterface,
) *Handler {
	authMiddleware := middleware.NewAuthMiddleware(useCases.User)

	return &Handler{
		userHandler:      NewUserHandler(useCases.User, logger),
		pvzHandler:       NewPVZHandler(useCases.PVZ, logger, metrics),
		receptionHandler: NewReceptionHandler(useCases.Reception, logger, metrics),
		productHandler:   NewProductHandler(useCases.Product, logger, metrics),
		authMiddleware:   authMiddleware,
		logger:           logger,
		metrics:          metrics,
	}
}

// Init инициализирует маршруты
func (h *Handler) Init(router *gin.Engine) {
	// Middleware для метрик
	router.Use(h.metricsMiddleware())

	// Группа API
	api := router.Group("/")
	{
		// Авторизация и регистрация
		api.POST("/dummyLogin", h.userHandler.DummyLogin)
		api.POST("/register", h.userHandler.Register)
		api.POST("/login", h.userHandler.Login)

		// Маршруты, требующие аутентификации
		authenticated := api.Group("/", h.authMiddleware.Authenticate())
		{
			// ПВЗ
			pvz := authenticated.Group("/pvz")
			{
				pvz.POST("/", h.authMiddleware.CheckRole(models.ModeratorRole), h.pvzHandler.Create)
				pvz.GET("/", h.pvzHandler.List)

				// Закрытие приемки и удаление товара
				pvz.POST("/:pvzId/close_last_reception", h.authMiddleware.CheckRole(models.EmployeeRole), h.receptionHandler.CloseLastReception)
				pvz.POST("/:pvzId/delete_last_product", h.authMiddleware.CheckRole(models.EmployeeRole), h.productHandler.DeleteLastFromReception)
			}

			// Приемки
			authenticated.POST("/receptions", h.authMiddleware.CheckRole(models.EmployeeRole), h.receptionHandler.Create)

			// Товары
			authenticated.POST("/products", h.authMiddleware.CheckRole(models.EmployeeRole), h.productHandler.Create)
		}
	}
}

// metricsMiddleware создает middleware для сбора метрик
func (h *Handler) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		// Метрики запроса
		h.metrics.IncRequestCount(
			c.Request.Method,
			c.FullPath(),
			strconv.Itoa(c.Writer.Status()),
		)

		// Метрики времени ответа
		h.metrics.ObserveRequestDuration(
			c.Request.Method,
			c.FullPath(),
			time.Since(startTime).Seconds(),
		)
	}
}
