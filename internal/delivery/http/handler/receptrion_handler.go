package handler

import (
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
)

// ReceptionHandler представляет обработчик для приемок
type ReceptionHandler struct {
	receptionUseCase usecase.ReceptionUseCase
	logger           logger.Logger
	metrics          metrics.MetricsInterface
}

func NewReceptionHandler(receptionUseCase usecase.ReceptionUseCase, logger logger.Logger, metrics metrics.MetricsInterface) *ReceptionHandler {
	return &ReceptionHandler{
		receptionUseCase: receptionUseCase,
		logger:           logger,
		metrics:          metrics,
	}
}

type createReceptionRequest struct {
	PVZID uuid.UUID `json:"pvzId" binding:"required"`
}

func (h *ReceptionHandler) Create(c *gin.Context) {
	var req createReceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	reception, err := h.receptionUseCase.Create(c.Request.Context(), req.PVZID)
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "pvz not found"})
			return
		}
		if err == errors.ErrOpenReceptionExists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "open reception already exists"})
			return
		}
		h.logger.Error("failed to create reception", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	h.metrics.IncReceptionCreated()

	c.JSON(http.StatusCreated, reception)
}

func (h *ReceptionHandler) CloseLastReception(c *gin.Context) {
	pvzIDStr := c.Param("pvzId")
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid pvz id"})
		return
	}

	reception, err := h.receptionUseCase.CloseLastReception(c.Request.Context(), pvzID)
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "pvz not found"})
			return
		}
		if err == errors.ErrOpenReceptionNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no open reception found"})
			return
		}
		h.logger.Error("failed to close reception", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, reception)
}
