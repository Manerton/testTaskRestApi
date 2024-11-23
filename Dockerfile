FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем исходный код приложения в контейнер
COPY . .

# Сборка приложения
RUN go build -o server ./cmd/main.go

# Второй этап: используем минимальный образ для запуска приложения
FROM debian:bookworm-slim

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Рабочая директория в контейнере
WORKDIR /app

# Копируем скомпилированный сервер и конфигурацию
COPY --from=builder /app/server /app/server
COPY config.env /app/config.env

# Указываем порт приложения
EXPOSE 8080

# Устанавливаем переменные окружения
ENV CONFIG_PATH=/app/config.env

# Запуск приложения
CMD ["/app/server"]