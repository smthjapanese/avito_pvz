package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	mock_usecase "github.com/smthjapanese/avito_pvz/internal/domain/usecase/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware(t *testing.T) {
	// Инициализация тестового окружения
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock_usecase.NewMockUserUseCase(ctrl)
	authMiddleware := NewAuthMiddleware(mockUserUseCase)

	t.Run("Authenticate - Success", func(t *testing.T) {
		// Подготовка тестовых данных
		token := "valid_token"
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
			Role:  models.ModeratorRole,
		}

		// Настройка мока
		mockUserUseCase.EXPECT().
			ValidateToken(gomock.Any(), token).
			Return(user, nil)

		// Создание тестового запроса
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// Выполнение middleware
		handler := authMiddleware.Authenticate()
		handler(c)

		// Проверка результатов
		assert.Equal(t, http.StatusOK, w.Code)
		userFromCtx, exists := c.Get(userCtx)
		assert.True(t, exists)
		assert.Equal(t, user, userFromCtx)
	})

	t.Run("Authenticate - Empty Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		handler := authMiddleware.Authenticate()
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Authenticate - Invalid Header Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "InvalidFormat")

		handler := authMiddleware.Authenticate()
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("CheckRole - Success", func(t *testing.T) {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
			Role:  models.ModeratorRole,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(userCtx, user)

		handler := authMiddleware.CheckRole(models.ModeratorRole)
		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CheckRole - Wrong Role", func(t *testing.T) {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
			Role:  models.EmployeeRole,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set(userCtx, user)

		handler := authMiddleware.CheckRole(models.ModeratorRole)
		handler(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("CheckRole - No User in Context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler := authMiddleware.CheckRole(models.ModeratorRole)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GetUser - Success", func(t *testing.T) {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
			Role:  models.ModeratorRole,
		}

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set(userCtx, user)

		result, err := GetUser(c)
		require.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("GetUser - No User in Context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		result, err := GetUser(c)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
