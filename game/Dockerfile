FROM golang:1.23-bullseye AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o game ./cmd

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

# Копируем исполняемый файл из builder
COPY --from=builder /app/game .

# Копируем .env файл
COPY .env .

# Открываем порт
EXPOSE 8083

# Запускаем приложение
CMD ["./game"] 