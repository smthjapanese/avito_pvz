package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/smthjapanese/avito_pvz/internal/app"
	"github.com/smthjapanese/avito_pvz/internal/config"
	"github.com/smthjapanese/avito_pvz/internal/delivery/grpc"
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

	// Создание и запуск приложения
	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Создание и запуск gRPC-сервера
	server := grpc.NewServer(application.GetPVZUseCase())
	if err := server.Start("3000"); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	// Обработка сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down servers...")

	server.Stop()
	log.Println("Server exited properly")
}
