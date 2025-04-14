package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smthjapanese/avito_pvz/internal/app"
	"github.com/smthjapanese/avito_pvz/internal/config"
)

func main() {
	// Загрузка конфигурации
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создание приложения
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Запуск приложения (HTTP, gRPC и metrics серверы)
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}

	// Обработка сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down servers...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем приложение
	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown app: %v", err)
	}

	log.Println("Server exited properly")
}
