FROM golang:1.23-bullseye AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway ./cmd/app

# Создаем минимальный образ для запуска
FROM alpine:3.18

# Устанавливаем CA сертификаты
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/api-gateway .

# Устанавливаем переменные окружения
ENV API_GATEWAY_PORT=8090
ENV AUTH_SERVICE_URL=http://auth-service:8080
ENV AUTH_GRPC_SERVICE_URL=auth-service:9090
ENV EDU_SERVICE_URL=http://edu-service:8081
ENV GAME_SERVICE_URL=http://game-service:8083
ENV GIN_MODE=debug

# Указываем порт
EXPOSE 8090

# Запускаем приложение
CMD ["./api-gateway"] 