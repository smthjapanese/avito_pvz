package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	t.Run("DefaultParams", func(t *testing.T) {
		params := DefaultParams()
		assert.Equal(t, uint32(64*1024), params.Memory)
		assert.Equal(t, uint32(3), params.Iterations)
		assert.Equal(t, uint8(2), params.Parallelism)
		assert.Equal(t, uint32(16), params.SaltLength)
		assert.Equal(t, uint32(32), params.KeyLength)
	})

	t.Run("Hash and Verify", func(t *testing.T) {
		password := "test-password-123"
		params := DefaultParams()

		// Тестируем хеширование
		hash, err := Hash(password, params)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Contains(t, hash, "$argon2id$")

		// Тестируем верификацию правильного пароля
		match, err := Verify(password, hash)
		require.NoError(t, err)
		assert.True(t, match)

		// Тестируем верификацию неправильного пароля
		match, err = Verify("wrong-password", hash)
		require.NoError(t, err)
		assert.False(t, match)
	})

	t.Run("Hash with Custom Params", func(t *testing.T) {
		password := "test-password-123"
		params := &Params{
			Memory:      32 * 1024,
			Iterations:  2,
			Parallelism: 1,
			SaltLength:  8,
			KeyLength:   16,
		}

		hash, err := Hash(password, params)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)

		match, err := Verify(password, hash)
		require.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("Invalid Hash Format", func(t *testing.T) {
		// Неверный формат хеша
		_, err := Verify("password", "invalid-hash-format")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidHash, err)

		// Неверная версия Argon2
		invalidVersionHash := "$argon2id$v=1$m=65536,t=3,p=2$c2FsdA$cGFzc3dvcmQ"
		_, err = Verify("password", invalidVersionHash)
		assert.Error(t, err)
		assert.Equal(t, ErrIncompatibleVersion, err)
	})

	t.Run("Empty Password", func(t *testing.T) {
		params := DefaultParams()
		hash, err := Hash("", params)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)

		match, err := Verify("", hash)
		require.NoError(t, err)
		assert.True(t, match)
	})
}
