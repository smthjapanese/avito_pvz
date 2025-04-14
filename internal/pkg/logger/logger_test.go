package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testLogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Caller  string `json:"caller"`
}

func TestLogger(t *testing.T) {
	t.Run("NewLogger with valid level", func(t *testing.T) {
		levels := []string{"debug", "info", "warn", "error", "fatal"}
		for _, level := range levels {
			logger, err := NewLogger(level)
			require.NoError(t, err)
			assert.NotNil(t, logger)
		}
	})

	t.Run("NewLogger with invalid level", func(t *testing.T) {
		logger, err := NewLogger("invalid")
		assert.Error(t, err)
		assert.Nil(t, logger)
	})

	t.Run("Logger levels", func(t *testing.T) {
		// Создаем буфер для записи логов
		var buf bytes.Buffer
		encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		})

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(&buf),
			zapcore.DebugLevel,
		)

		logger := zap.New(core, zap.AddCaller())
		defer logger.Sync()

		// Тестируем разные уровни логирования
		testCases := []struct {
			level   string
			message string
			fields  []zapcore.Field
		}{
			{"debug", "debug message", []zapcore.Field{zap.String("key", "value")}},
			{"info", "info message", []zapcore.Field{zap.Int("count", 42)}},
			{"warn", "warn message", []zapcore.Field{zap.Bool("flag", true)}},
			{"error", "error message", []zapcore.Field{zap.Error(nil)}},
		}

		for _, tc := range testCases {
			buf.Reset()
			switch tc.level {
			case "debug":
				logger.Debug(tc.message, tc.fields...)
			case "info":
				logger.Info(tc.message, tc.fields...)
			case "warn":
				logger.Warn(tc.message, tc.fields...)
			case "error":
				logger.Error(tc.message, tc.fields...)
			}

			// Проверяем, что лог был записан
			assert.NotEmpty(t, buf.String())

			// Проверяем структуру JSON
			var entry testLogEntry
			err := json.Unmarshal(buf.Bytes(), &entry)
			require.NoError(t, err)
			assert.Equal(t, tc.level, entry.Level)
			assert.Equal(t, tc.message, entry.Message)
			assert.NotEmpty(t, entry.Caller)
		}
	})

	t.Run("Logger with different log levels", func(t *testing.T) {
		// Создаем буфер для записи логов
		var buf bytes.Buffer
		encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		})

		// Создаем логгер с уровнем INFO
		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(&buf),
			zapcore.InfoLevel,
		)

		logger := zap.New(core, zap.AddCaller())
		defer logger.Sync()

		// Отправляем сообщения разных уровней
		logger.Debug("debug message") // Не должно быть записано
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		// Проверяем, что debug сообщение не было записано
		assert.NotContains(t, buf.String(), "debug message")
		// Проверяем, что остальные сообщения были записаны
		assert.Contains(t, buf.String(), "info message")
		assert.Contains(t, buf.String(), "warn message")
		assert.Contains(t, buf.String(), "error message")
	})
}
