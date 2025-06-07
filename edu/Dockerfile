FROM golang:1.23-bullseye AS builder

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o edu-service ./cmd

# Этап запуска
FROM ubuntu:22.04

# Устанавливаем необходимые пакеты
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/edu-service /app/edu-service

# Копируем конфигурационные файлы
COPY --from=builder /app/.env /app/.env

# Открываем порты (HTTP и gRPC)
EXPOSE 8081 9090

# Запуск сервиса
CMD ["./edu-service"] 