package handler

import (
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCase
	logger         logger.Logger
	metrics        metrics.MetricsInterface
}

func NewProductHandler(productUseCase usecase.ProductUseCase, logger logger.Logger, metrics metrics.MetricsInterface) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		logger:         logger,
		metrics:        metrics,
	}
}

type createProductRequest struct {
	Type  models.ProductType `json:"type" binding:"required"`
	PVZID uuid.UUID          `json:"pvzId" binding:"required"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	product, err := h.productUseCase.Create(c.Request.Context(), req.Type, req.PVZID)
	if err != nil {
		if err == errors.ErrInvalidProductType {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid product type"})
			return
		}
		if errors.IsNotFound(err) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "pvz not found"})
			return
		}
		if err == errors.ErrOpenReceptionNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no open reception found"})
			return
		}
		h.logger.Error("failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	h.metrics.IncProductAdded()

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) DeleteLastFromReception(c *gin.Context) {
	pvzIDStr := c.Param("pvzId")
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid pvz id"})
		return
	}

	err = h.productUseCase.DeleteLastFromReception(c.Request.Context(), pvzID)
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "pvz not found"})
			return
		}
		if err == errors.ErrOpenReceptionNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no open reception found"})
			return
		}
		if err == errors.ErrNoProductsToDelete {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no products to delete"})
			return
		}
		h.logger.Error("failed to delete product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
