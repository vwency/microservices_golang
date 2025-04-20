# --- Build stage ---
    FROM golang:1.24.2-alpine AS builder
    WORKDIR /app
    
    # Загрузка зависимостей
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Копируем исходники и билдим
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/database_service cmd/database_service/main.go
    
    # --- Final stage ---
    FROM alpine:3.19
    # Устанавливаем зависимости
    RUN apk add --no-cache ca-certificates postgresql-client bash
    
    # Добавляем wait-for-it скрипт
    ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
    RUN chmod +x /wait-for-it.sh
    
    # Копируем бинарник
    COPY --from=builder /bin/database_service /bin/database_service
    
    # Копируем конфиги
    COPY ./config /app/config
    
    # Runtime переменные
    ENV ENV=dev \
        PORT=50052 \
        LOG_LEVEL=debug \
        DATABASE_URL="postgres://postgres:password@postgres:5432/pure?sslmode=disable"
    
    WORKDIR /app
    
    # Запуск
    ENTRYPOINT ["/bin/sh", "-c", "/wait-for-it.sh postgres:5432 -t 6 -- /bin/database_service"]
