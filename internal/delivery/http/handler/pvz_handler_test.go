package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	mock_usecase "github.com/smthjapanese/avito_pvz/internal/domain/usecase/mock"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPVZHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZUseCase := mock_usecase.NewMockPVZUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()
	handler := NewPVZHandler(mockPVZUseCase, mockLogger, mockMetrics)

	req := createPVZRequest{
		City: models.CityMoscow,
	}
	reqBody, _ := json.Marshal(req)

	pvz := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             req.City,
		CreatedAt:        time.Now(),
	}
	mockPVZUseCase.EXPECT().Create(gomock.Any(), req.City).Return(pvz, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.PVZ
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, pvz.ID, response.ID)
	assert.Equal(t, pvz.City, response.City)
}

func TestPVZHandler_Create_InvalidCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZUseCase := mock_usecase.NewMockPVZUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics() // Используем мок метрик
	handler := NewPVZHandler(mockPVZUseCase, mockLogger, mockMetrics)

	req := createPVZRequest{
		City: "InvalidCity",
	}
	reqBody, _ := json.Marshal(req)

	mockPVZUseCase.EXPECT().Create(gomock.Any(), models.City(req.City)).Return(nil, errors.ErrInvalidCity)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/pvz", handler.Create)

	c.Request, _ = http.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid city")
}

func TestPVZHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZUseCase := mock_usecase.NewMockPVZUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics() // Используем мок метрик
	handler := NewPVZHandler(mockPVZUseCase, mockLogger, mockMetrics)

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	page := 1
	limit := 10

	pvz1 := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now().Add(-12 * time.Hour),
		City:             models.CityMoscow,
		CreatedAt:        time.Now().Add(-12 * time.Hour),
	}

	reception1 := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now().Add(-10 * time.Hour),
		PVZID:     pvz1.ID,
		Status:    models.ReceptionStatusClose,
		CreatedAt: time.Now().Add(-10 * time.Hour),
	}

	product1 := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now().Add(-9 * time.Hour),
		Type:        models.ProductTypeElectronics,
		ReceptionID: reception1.ID,
		CreatedAt:   time.Now().Add(-9 * time.Hour),
	}

	receptionWithProducts := &usecase.ReceptionWithProducts{
		Reception: reception1,
		Products:  []*models.Product{product1},
	}

	pvzWithReceptions := []*usecase.PVZWithReceptions{
		{
			PVZ:        pvz1,
			Receptions: []*usecase.ReceptionWithProducts{receptionWithProducts},
		},
	}

	mockPVZUseCase.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), page, limit).
		DoAndReturn(func(_ interface{}, startDateParam, endDateParam *time.Time, pageParam, limitParam int) ([]*usecase.PVZWithReceptions, error) {
			// Проверяем, что параметры соответствуют ожидаемым
			assert.NotNil(t, startDateParam)
			assert.NotNil(t, endDateParam)
			assert.True(t, startDateParam.Equal(startDate) || startDateParam.Sub(startDate) < time.Second)
			assert.True(t, endDateParam.Equal(endDate) || endDateParam.Sub(endDate) < time.Second)
			assert.Equal(t, page, pageParam)
			assert.Equal(t, limit, limitParam)

			return pvzWithReceptions, nil
		})

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/pvz", handler.List)

	startDateStr := url.QueryEscape(startDate.Format(time.RFC3339))
	endDateStr := url.QueryEscape(endDate.Format(time.RFC3339))
	urlPath := fmt.Sprintf("/pvz?startDate=%s&endDate=%s&page=%d&limit=%d", startDateStr, endDateStr, page, limit)

	c.Request, _ = http.NewRequest(http.MethodGet, urlPath, nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)

	// Если ответ успешный, проверяем его содержимое
	if w.Code == http.StatusOK {
		var response []*usecase.PVZWithReceptions
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, pvz1.ID, response[0].PVZ.ID)
		assert.Equal(t, pvz1.City, response[0].PVZ.City)
		assert.Len(t, response[0].Receptions, 1)
		assert.Equal(t, reception1.ID, response[0].Receptions[0].Reception.ID)
		assert.Len(t, response[0].Receptions[0].Products, 1)
		assert.Equal(t, product1.ID, response[0].Receptions[0].Products[0].ID)
	} else {
		// Если тест не проходит, выводим ответ для отладки
		t.Logf("Response body: %s", w.Body.String())
	}
}

func TestPVZHandler_List_InvalidDateFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZUseCase := mock_usecase.NewMockPVZUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics() // Используем мок метрик
	handler := NewPVZHandler(mockPVZUseCase, mockLogger, mockMetrics)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/pvz", handler.List)

	c.Request, _ = http.NewRequest(http.MethodGet, "/pvz?startDate=invalid-date", nil)

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid start date format")
}
