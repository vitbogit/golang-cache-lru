# Указываем базовый образ с Go
FROM golang:1.23-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o cacheapp ./cmd/lru/lru.go

# Используем минимальный образ для запуска приложения
FROM alpine:latest

# Install any runtime dependencies (optional, usually none for static Go binaries)
RUN apk add --no-cache ca-certificates

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/cacheapp /cacheapp

# Копируем файлы конфигурации
COPY --from=builder /app/configs ./configs

# Копируем файл .env (если используется)
COPY --from=builder /app/.env .env

# Указываем порт, который будет использоваться приложением
EXPOSE 8080

# Команда запуска приложения
CMD ["./cacheapp"]
