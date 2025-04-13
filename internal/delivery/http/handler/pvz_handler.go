package handler

import (
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
)

// PVZHandler представляет обработчик для ПВЗ
type PVZHandler struct {
	pvzUseCase usecase.PVZUseCase
	logger     logger.Logger
	metrics    metrics.MetricsInterface
}

func NewPVZHandler(pvzUseCase usecase.PVZUseCase, logger logger.Logger, metrics metrics.MetricsInterface) *PVZHandler {
	return &PVZHandler{
		pvzUseCase: pvzUseCase,
		logger:     logger,
		metrics:    metrics,
	}
}

type createPVZRequest struct {
	City models.City `json:"city" binding:"required"`
}

func (h *PVZHandler) Create(c *gin.Context) {
	var req createPVZRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	pvz, err := h.pvzUseCase.Create(c.Request.Context(), req.City)
	if err != nil {
		if err == errors.ErrInvalidCity {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid city"})
			return
		}
		h.logger.Error("failed to create pvz", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	h.metrics.IncPVZCreated()

	c.JSON(http.StatusCreated, pvz)
}

type listPVZRequest struct {
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=10" binding:"min=1,max=30"`
}

func (h *PVZHandler) List(c *gin.Context) {
	var req listPVZRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var startDate, endDate *time.Time
	if req.StartDate != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid start date format"})
			return
		}
		startDate = &parsedStartDate
	}

	if req.EndDate != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid end date format"})
			return
		}
		endDate = &parsedEndDate
	}

	pvzs, err := h.pvzUseCase.List(c.Request.Context(), startDate, endDate, req.Page, req.Limit)
	if err != nil {
		h.logger.Error("failed to get pvz list", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, pvzs)
}
