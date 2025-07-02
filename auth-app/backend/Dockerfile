FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main .

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированный бинарник
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]