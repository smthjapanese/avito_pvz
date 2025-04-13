package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase/mock"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()
	receptionID := uuid.New()
	req := createProductRequest{
		Type:  models.ProductTypeElectronics,
		PVZID: pvzID,
	}
	reqBody, _ := json.Marshal(req)

	product := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        req.Type,
		ReceptionID: receptionID,
		CreatedAt:   time.Now(),
	}
	mockProductUseCase.EXPECT().Create(gomock.Any(), req.Type, req.PVZID).Return(product, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/products", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, product.ID, response.ID)
	assert.Equal(t, product.Type, response.Type)
	assert.Equal(t, product.ReceptionID, response.ReceptionID)
}

func TestProductHandler_Create_InvalidProductType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()
	req := createProductRequest{
		Type:  "invalid-type",
		PVZID: pvzID,
	}
	reqBody, _ := json.Marshal(req)

	mockProductUseCase.EXPECT().Create(gomock.Any(), models.ProductType(req.Type), req.PVZID).Return(nil, errors.ErrInvalidProductType)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/products", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid product type")
}

func TestProductHandler_Create_NoOpenReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()
	req := createProductRequest{
		Type:  models.ProductTypeElectronics,
		PVZID: pvzID,
	}
	reqBody, _ := json.Marshal(req)

	mockProductUseCase.EXPECT().Create(gomock.Any(), req.Type, req.PVZID).Return(nil, errors.ErrPVZNotFound)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/products", handler.Create)
	c.Request, _ = http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "pvz not found")
}

func TestProductHandler_DeleteLastFromReception(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()

	mockProductUseCase.EXPECT().DeleteLastFromReception(gomock.Any(), pvzID).Return(nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/delete_last_product", handler.DeleteLastFromReception)

	c.Params = []gin.Param{
		{Key: "pvzId", Value: pvzID.String()},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "product deleted")
}

func TestProductHandler_DeleteLastFromReception_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/delete_last_product", handler.DeleteLastFromReception)

	c.Params = []gin.Param{
		{Key: "pvzId", Value: "invalid-uuid"},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/invalid-uuid/delete_last_product", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid pvz id")
}

func TestProductHandler_DeleteLastFromReception_NoProductsToDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()

	mockProductUseCase.EXPECT().DeleteLastFromReception(gomock.Any(), pvzID).Return(errors.ErrNoProductsToDelete)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/delete_last_product", handler.DeleteLastFromReception)

	c.Params = []gin.Param{
		{Key: "pvzId", Value: pvzID.String()},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "no products to delete")
}

func TestProductHandler_DeleteLastFromReception_NoOpenReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()

	customErr := errors.ErrOpenReceptionNotFound

	mockProductUseCase.EXPECT().DeleteLastFromReception(gomock.Any(), pvzID).Return(customErr)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/delete_last_product", handler.DeleteLastFromReception)
	c.Params = []gin.Param{
		{Key: "pvzId", Value: pvzID.String()},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "pvz not found")
}

func TestProductHandler_DeleteLastFromReception_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()

	mockProductUseCase.EXPECT().DeleteLastFromReception(gomock.Any(), pvzID).Return(errors.ErrPVZNotFound)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/delete_last_product", handler.DeleteLastFromReception)

	c.Params = []gin.Param{
		{Key: "pvzId", Value: pvzID.String()},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "pvz not found")
}

func TestProductHandler_Create_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()
	req := createProductRequest{
		Type:  models.ProductTypeElectronics,
		PVZID: pvzID,
	}
	reqBody, _ := json.Marshal(req)

	mockProductUseCase.EXPECT().Create(gomock.Any(), req.Type, req.PVZID).Return(nil, errors.ErrPVZNotFound)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/products", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "pvz not found")
}

func TestProductHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	reqBody := []byte(`{}`)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/products", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "message")
}

func TestProductHandler_Create_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()
	req := createProductRequest{
		Type:  models.ProductTypeElectronics,
		PVZID: pvzID,
	}
	reqBody, _ := json.Marshal(req)

	mockProductUseCase.EXPECT().Create(gomock.Any(), req.Type, req.PVZID).Return(nil, errors.ErrInternal)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/products", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
}

func TestProductHandler_DeleteLastFromReception_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductUseCase := mock.NewMockProductUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewProductHandler(mockProductUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()

	mockProductUseCase.EXPECT().DeleteLastFromReception(gomock.Any(), pvzID).Return(errors.ErrInternal)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/delete_last_product", handler.DeleteLastFromReception)

	c.Params = []gin.Param{
		{Key: "pvzId", Value: pvzID.String()},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
}
