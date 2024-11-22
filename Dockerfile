# Используем официальный образ Go
FROM golang:alpine AS builder

ENV PORT="80"
ENV DB_URL="babich:babich@amvera-babich-run-tbank-db/go_api_medium?parseTime=true"

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

ENV DB_URL="babich:babich@amvera-babich-run-tbank-db:3306/go_api_medium?parseTime=true"

# Команда для запуска сервера
CMD ["./server"]