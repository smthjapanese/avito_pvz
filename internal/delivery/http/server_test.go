package http

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Run("NewServer", func(t *testing.T) {
		handler := http.NewServeMux()
		server := NewServer("8080", handler)

		assert.NotNil(t, server)
		assert.NotNil(t, server.httpServer)
		assert.Equal(t, ":8080", server.httpServer.Addr)
		assert.Equal(t, 10*time.Second, server.httpServer.ReadTimeout)
		assert.Equal(t, 10*time.Second, server.httpServer.WriteTimeout)
		assert.Equal(t, 1<<20, server.httpServer.MaxHeaderBytes)
	})

	t.Run("Run and Shutdown", func(t *testing.T) {
		// Создаем тестовый сервер
		handler := http.NewServeMux()
		server := NewServer("0", handler) // порт 0 для автоматического выбора свободного порта

		// Запускаем сервер в отдельной горутине
		errChan := make(chan error, 1)
		go func() {
			errChan <- server.Run()
		}()

		// Даем серверу время на запуск
		time.Sleep(100 * time.Millisecond)

		// Проверяем, что сервер запущен
		select {
		case err := <-errChan:
			require.NoError(t, err)
		default:
			// Сервер все еще работает, что хорошо
		}

		// Останавливаем сервер
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		require.NoError(t, err)
	})

	t.Run("Run Error", func(t *testing.T) {
		// Создаем сервер с невалидным портом
		handler := http.NewServeMux()
		server := NewServer("invalid-port", handler)

		// Запускаем сервер
		err := server.Run()
		require.Error(t, err)
	})

	t.Run("Shutdown Timeout", func(t *testing.T) {
		// Создаем тестовый сервер с обработчиком, который блокирует завершение
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
		})
		server := NewServer("0", handler)

		// Запускаем сервер в отдельной горутине
		go func() {
			_ = server.Run()
		}()

		// Даем серверу время на запуск
		time.Sleep(100 * time.Millisecond)

		// Создаем контекст с очень маленьким таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()

		// Делаем запрос к серверу, чтобы заблокировать его
		go http.Get("http://localhost:" + server.httpServer.Addr)

		// Пытаемся остановить сервер
		err := server.Shutdown(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})

	t.Run("Shutdown Error", func(t *testing.T) {
		// Создаем тестовый сервер с обработчиком, который блокирует завершение
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
		})
		server := NewServer("0", handler)

		// Запускаем сервер в отдельной горутине
		go func() {
			_ = server.Run()
		}()

		// Даем серверу время на запуск
		time.Sleep(100 * time.Millisecond)

		// Делаем запрос к серверу, чтобы заблокировать его
		go http.Get("http://localhost:" + server.httpServer.Addr)

		// Создаем контекст с отменой
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Отменяем контекст сразу

		// Пытаемся остановить сервер
		err := server.Shutdown(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})
}
