# Этап 1: Сборка приложения
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код приложения
COPY . .

# Сборка приложения в релизном режиме
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server main.go

# Этап 2: Минимальный образ для запуска
FROM scratch

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/server /server

# Указываем переменные окружения (если нужно)
ENV GIN_MODE=release

# Указываем порт, который будет прослушивать приложение
EXPOSE 8080

# Команда для запуска приложения
ENTRYPOINT ["/server"]
