package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mock_usecase "github.com/smthjapanese/avito_pvz/internal/domain/usecase/mock"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
	"github.com/smthjapanese/avito_pvz/internal/pkg/metrics"
	"github.com/smthjapanese/avito_pvz/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZUseCase := mock_usecase.NewMockPVZUseCase(ctrl)
	mockReceptionUseCase := mock_usecase.NewMockReceptionUseCase(ctrl)
	mockProductUseCase := mock_usecase.NewMockProductUseCase(ctrl)
	mockUserUseCase := mock_usecase.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()

	useCases := &usecase.UseCases{
		PVZ:       mockPVZUseCase,
		Reception: mockReceptionUseCase,
		Product:   mockProductUseCase,
		User:      mockUserUseCase,
	}

	handler := NewHandler(useCases, mockLogger, mockMetrics)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.pvzHandler)
	assert.NotNil(t, handler.receptionHandler)
	assert.NotNil(t, handler.productHandler)
	assert.NotNil(t, handler.userHandler)
}

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZUseCase := mock_usecase.NewMockPVZUseCase(ctrl)
	mockReceptionUseCase := mock_usecase.NewMockReceptionUseCase(ctrl)
	mockProductUseCase := mock_usecase.NewMockProductUseCase(ctrl)
	mockUserUseCase := mock_usecase.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()

	useCases := &usecase.UseCases{
		PVZ:       mockPVZUseCase,
		Reception: mockReceptionUseCase,
		Product:   mockProductUseCase,
		User:      mockUserUseCase,
	}

	handler := NewHandler(useCases, mockLogger, mockMetrics)

	router := gin.New()
	handler.Init(router)

	// Проверяем, что все маршруты зарегистрированы
	routes := router.Routes()
	assert.NotEmpty(t, routes)

	// Проверяем наличие основных маршрутов
	expectedRoutes := map[string]bool{
		"POST /register":                        false,
		"POST /login":                           false,
		"POST /dummyLogin":                      false,
		"POST /pvz/":                            false,
		"GET /pvz/":                             false,
		"POST /receptions":                      false,
		"POST /pvz/:pvzId/close_last_reception": false,
		"POST /products":                        false,
		"POST /pvz/:pvzId/delete_last_product":  false,
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, exists := expectedRoutes[key]; exists {
			expectedRoutes[key] = true
		}
	}

	for route, found := range expectedRoutes {
		assert.True(t, found, "Route %s not found", route)
	}
}

func TestMetricsMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZUseCase := mock_usecase.NewMockPVZUseCase(ctrl)
	mockReceptionUseCase := mock_usecase.NewMockReceptionUseCase(ctrl)
	mockProductUseCase := mock_usecase.NewMockProductUseCase(ctrl)
	mockUserUseCase := mock_usecase.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()

	useCases := &usecase.UseCases{
		PVZ:       mockPVZUseCase,
		Reception: mockReceptionUseCase,
		Product:   mockProductUseCase,
		User:      mockUserUseCase,
	}

	handler := NewHandler(useCases, mockLogger, mockMetrics)

	// Создаем тестовый запрос
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Применяем middleware
	middleware := handler.metricsMiddleware()
	middleware(c)

	// Проверяем, что запрос обработан
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMetricsMiddlewareError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZUseCase := mock_usecase.NewMockPVZUseCase(ctrl)
	mockReceptionUseCase := mock_usecase.NewMockReceptionUseCase(ctrl)
	mockProductUseCase := mock_usecase.NewMockProductUseCase(ctrl)
	mockUserUseCase := mock_usecase.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	mockMetrics := metrics.NewMockMetrics()

	useCases := &usecase.UseCases{
		PVZ:       mockPVZUseCase,
		Reception: mockReceptionUseCase,
		Product:   mockProductUseCase,
		User:      mockUserUseCase,
	}

	handler := NewHandler(useCases, mockLogger, mockMetrics)

	// Создаем тестовый запрос с ошибкой
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Добавляем обработчик с ошибкой
	r.GET("/test", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusInternalServerError)
	})

	// Применяем middleware и выполняем запрос
	r.Use(handler.metricsMiddleware())
	r.ServeHTTP(w, c.Request)

	// Проверяем, что запрос обработан с ошибкой
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
