package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smthjapanese/avito_pvz/internal/delivery/http/middleware"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
	"github.com/smthjapanese/avito_pvz/internal/usecase"
)

type Handler struct {
	userHandler      *UserHandler
	pvzHandler       *PVZHandler
	receptionHandler *ReceptionHandler
	productHandler   *ProductHandler
	authMiddleware   *middleware.AuthMiddleware
	logger           logger.Logger
	metrics          metrics.MetricsInterface
}

func NewHandler(useCases *usecase.UseCases, logger logger.Logger, metrics metrics.MetricsInterface) *Handler {
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

func (h *Handler) Init(router *gin.Engine) {
	router.Use(h.metricsMiddleware())

	api := router.Group("/")
	{
		// Авторизация и регистрация
		api.POST("/dummyLogin", h.userHandler.DummyLogin)
		api.POST("/register", h.userHandler.Register)
		api.POST("/login", h.userHandler.Login)

		authenticated := api.Group("/", h.authMiddleware.Authenticate())
		{

			pvz := authenticated.Group("/pvz")
			{
				pvz.POST("/", h.authMiddleware.CheckRole(models.ModeratorRole), h.pvzHandler.Create)
				pvz.GET("/", h.pvzHandler.List)

				pvz.POST("/:pvzId/close_last_reception", h.authMiddleware.CheckRole(models.EmployeeRole), h.receptionHandler.CloseLastReception)
				pvz.POST("/:pvzId/delete_last_product", h.authMiddleware.CheckRole(models.EmployeeRole), h.productHandler.DeleteLastFromReception)
			}

			authenticated.POST("/receptions", h.authMiddleware.CheckRole(models.EmployeeRole), h.receptionHandler.Create)

			authenticated.POST("/products", h.authMiddleware.CheckRole(models.EmployeeRole), h.productHandler.Create)
		}
	}
}

func (h *Handler) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		h.metrics.IncRequestCount(
			c.Request.Method,
			c.FullPath(),
			strconv.Itoa(c.Writer.Status()),
		)

		h.metrics.ObserveRequestDuration(
			c.Request.Method,
			c.FullPath(),
			time.Since(startTime).Seconds(),
		)
	}
}
