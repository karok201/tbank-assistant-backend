# Используем официальный образ Go
FROM golang:alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы для установки зависимостей
COPY go.mod go.sum .env ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Собираем приложение
RUN go build -o server .

# Минимизируем размер конечного образа
FROM alpine:latest

# Добавляем необходимые зависимости для запуска
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем собранный сервер из предыдущего этапа
COPY --from=builder /app/server .

# Указываем порт, на котором работает сервер
EXPOSE 8080

ENV DB_URL="babich:babich@tcp(127.0.0.1:3306)/tbank-db?charset=utf8mb4&parseTime=True&loc=Local"
# Команда для запуска сервера
CMD ["./server"]