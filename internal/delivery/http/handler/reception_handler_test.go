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

func TestReceptionHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionUseCase := mock.NewMockReceptionUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMetrics()
	handler := NewReceptionHandler(mockReceptionUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()
	req := createReceptionRequest{
		PVZID: pvzID,
	}
	reqBody, _ := json.Marshal(req)

	reception := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}
	mockReceptionUseCase.EXPECT().Create(gomock.Any(), pvzID).Return(reception, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/receptions", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Reception
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, reception.ID, response.ID)
	assert.Equal(t, reception.PVZID, response.PVZID)
	assert.Equal(t, reception.Status, response.Status)
}

func TestReceptionHandler_Create_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionUseCase := mock.NewMockReceptionUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMetrics()
	handler := NewReceptionHandler(mockReceptionUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()
	req := createReceptionRequest{
		PVZID: pvzID,
	}
	reqBody, _ := json.Marshal(req)

	mockReceptionUseCase.EXPECT().Create(gomock.Any(), pvzID).Return(nil, errors.ErrPVZNotFound)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/receptions", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "pvz not found")
}

func TestReceptionHandler_Create_OpenReceptionExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionUseCase := mock.NewMockReceptionUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMetrics()
	handler := NewReceptionHandler(mockReceptionUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()
	req := createReceptionRequest{
		PVZID: pvzID,
	}
	reqBody, _ := json.Marshal(req)

	mockReceptionUseCase.EXPECT().Create(gomock.Any(), pvzID).Return(nil, errors.ErrOpenReceptionExists)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/receptions", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "open reception already exists")
}

func TestReceptionHandler_CloseLastReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionUseCase := mock.NewMockReceptionUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMetrics()
	handler := NewReceptionHandler(mockReceptionUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()

	reception := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusClose,
		CreatedAt: time.Now(),
	}
	mockReceptionUseCase.EXPECT().CloseLastReception(gomock.Any(), pvzID).Return(reception, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/close_last_reception", handler.CloseLastReception)

	c.Params = []gin.Param{
		{Key: "pvzId", Value: pvzID.String()},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/close_last_reception", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Reception
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, reception.ID, response.ID)
	assert.Equal(t, reception.Status, response.Status)
}

func TestReceptionHandler_CloseLastReception_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionUseCase := mock.NewMockReceptionUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMetrics()
	handler := NewReceptionHandler(mockReceptionUseCase, mockLogger, mockMetrics)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/close_last_reception", handler.CloseLastReception)

	c.Params = []gin.Param{
		{Key: "pvzId", Value: "invalid-uuid"},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/invalid-uuid/close_last_reception", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid pvz id")
}

func TestReceptionHandler_CloseLastReception_NoOpenReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockReceptionUseCase := mock.NewMockReceptionUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics() // Используем мок метрик
	handler := NewReceptionHandler(mockReceptionUseCase, mockLogger, mockMetrics)

	pvzID := uuid.New()

	mockReceptionUseCase.EXPECT().CloseLastReception(gomock.Any(), pvzID).Return(nil, errors.ErrPVZNotFound)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz/:pvzId/close_last_reception", handler.CloseLastReception)
	c.Params = []gin.Param{
		{Key: "pvzId", Value: pvzID.String()},
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/close_last_reception", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "pvz not found")
}
