package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager(t *testing.T) {
	signingKey := "test-key"
	expiration := 1 * time.Hour
	manager := NewManager(signingKey, expiration)

	t.Run("GenerateToken", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"
		role := models.ModeratorRole

		token, err := manager.GenerateToken(userID, email, role)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Проверяем, что токен можно распарсить
		claims, err := manager.ParseToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID.String(), claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("GenerateDummyToken", func(t *testing.T) {
		role := models.EmployeeRole

		token, err := manager.GenerateDummyToken(role)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Проверяем, что токен можно распарсить
		claims, err := manager.ParseToken(token)
		require.NoError(t, err)
		assert.Equal(t, role, claims.Role)
		assert.Contains(t, claims.Email, "dummy_")
	})

	t.Run("ParseToken - Invalid Token", func(t *testing.T) {
		// Неверный формат токена
		_, err := manager.ParseToken("invalid.token.format")
		assert.Error(t, err)

		// Неверная подпись
		otherManager := NewManager("different-key", expiration)
		token, _ := otherManager.GenerateToken(uuid.New(), "test@example.com", models.ModeratorRole)
		_, err = manager.ParseToken(token)
		assert.Error(t, err)
	})

	t.Run("ParseToken - Expired Token", func(t *testing.T) {
		// Создаем менеджер с очень коротким сроком жизни токена
		shortExpiration := 1 * time.Millisecond
		shortManager := NewManager(signingKey, shortExpiration)

		token, err := shortManager.GenerateToken(uuid.New(), "test@example.com", models.ModeratorRole)
		require.NoError(t, err)

		// Ждем, пока токен истечет
		time.Sleep(2 * time.Millisecond)

		_, err = shortManager.ParseToken(token)
		assert.Error(t, err)
	})
}
