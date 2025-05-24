# Этап сборки
FROM golang:1.23-bullseye AS builder

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd

# Этап запуска
FROM ubuntu:22.04

# Устанавливаем необходимые пакеты (например, для отправки email может понадобиться ca-certificates)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/auth-service /app/auth-service

# Копируем конфигурационные файлы
COPY --from=builder /app/.env /app/.env

# Открываем порты (HTTP и gRPC)
EXPOSE 8080 9090

# Запуск сервиса
CMD ["./auth-service"]