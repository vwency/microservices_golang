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
    RUN apk add --no-cache ca-certificates postgresql-client
    
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
    
    # Healthcheck
    HEALTHCHECK --interval=30s --timeout=3s CMD pg_isready -d "$DATABASE_URL" || exit 1
    
    # Запуск
    ENTRYPOINT ["/bin/database_service"]
    