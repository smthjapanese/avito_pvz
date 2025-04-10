.PHONY: build run test lint clean docker-build docker-up docker-down

# Сборка приложения
build:
	go build -o ./bin/app ./cmd/app/main.go

# Запуск приложения
run:
	go run ./cmd/app/main.go

# Запуск тестов
test:
	go test -v ./...

# Проверка кода линтером
lint:
	golangci-lint run ./...

# Очистка артефактов сборки
clean:
	rm -rf ./bin

# Сборка Docker образа
docker-build:
	docker-compose build

# Запуск контейнеров
docker-up:
	docker-compose up -d

# Остановка контейнеров
docker-down:
	docker-compose down
