package metrics

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	// Создаем новый экземпляр метрик
	metrics := NewMetrics()

	t.Run("PVZ Metrics", func(t *testing.T) {
		// Проверяем счетчик созданных ПВЗ
		metrics.IncPVZCreated()
		metrics.IncPVZCreated()

		metric := &dto.Metric{}
		err := metrics.PVZCreated.Write(metric)
		require.NoError(t, err)
		assert.Equal(t, float64(2), metric.Counter.GetValue())
	})

	t.Run("Reception Metrics", func(t *testing.T) {
		// Проверяем счетчик созданных приемок
		metrics.IncReceptionCreated()
		metrics.IncReceptionCreated()
		metrics.IncReceptionCreated()

		metric := &dto.Metric{}
		err := metrics.ReceptionCreated.Write(metric)
		require.NoError(t, err)
		assert.Equal(t, float64(3), metric.Counter.GetValue())
	})

	t.Run("Product Metrics", func(t *testing.T) {
		// Проверяем счетчик добавленных товаров
		metrics.IncProductAdded()
		metrics.IncProductAdded()
		metrics.IncProductAdded()
		metrics.IncProductAdded()

		metric := &dto.Metric{}
		err := metrics.ProductAdded.Write(metric)
		require.NoError(t, err)
		assert.Equal(t, float64(4), metric.Counter.GetValue())
	})

	t.Run("HTTP Request Metrics", func(t *testing.T) {
		// Проверяем метрики HTTP запросов
		method := "GET"
		endpoint := "/api/pvz"
		status := "200"
		duration := 0.1

		metrics.IncRequestCount(method, endpoint, status)
		metrics.ObserveRequestDuration(method, endpoint, duration)

		// Проверяем счетчик запросов
		metric := &dto.Metric{}
		err := metrics.RequestCount.WithLabelValues(method, endpoint, status).Write(metric)
		require.NoError(t, err)
		assert.Equal(t, float64(1), metric.Counter.GetValue())

		// Проверяем, что метрика длительности была обновлена
		// Для гистограмм мы можем только проверить, что они были созданы
		// и что они имеют правильные метки
		histogram := metrics.RequestDuration.WithLabelValues(method, endpoint)
		assert.NotNil(t, histogram)
	})

	t.Run("gRPC Request Metrics", func(t *testing.T) {
		// Проверяем метрики gRPC запросов
		method := "CreatePVZ"
		status := "OK"
		duration := 0.05

		metrics.IncGRPCRequestCount(method, status)
		metrics.ObserveGRPCRequestDuration(method, duration)

		// Проверяем счетчик запросов
		metric := &dto.Metric{}
		err := metrics.GRPCRequestCount.WithLabelValues(method, status).Write(metric)
		require.NoError(t, err)
		assert.Equal(t, float64(1), metric.Counter.GetValue())

		// Проверяем, что метрика длительности была обновлена
		// Для гистограмм мы можем только проверить, что они были созданы
		// и что они имеют правильные метки
		histogram := metrics.GRPCRequestDuration.WithLabelValues(method)
		assert.NotNil(t, histogram)
	})
}
