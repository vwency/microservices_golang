# --- Build stage ---
    FROM golang:1.24.2-alpine AS builder
    WORKDIR /app
    
    # Загрузка зависимостей
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Копируем исходники и билдим
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/gateway cmd/gateway/main.go
    
    # --- Final stage ---
    FROM alpine:3.19
    # Устанавливаем зависимости
    RUN apk add --no-cache ca-certificates
    
    # Копируем бинарник
    COPY --from=builder /bin/gateway /bin/gateway
    
    # Копируем конфиги
    COPY ./config/gateway /app/config
    
    # Runtime переменные
    ENV ENV=dev \
        PORT=50058 \
        LOG_LEVEL=debug \
        AUTH_SERVICE_URL="localhost:50058" \
        DATABASE_URL="postgres://postgres:password@postgres:5438/pure?sslmode=disable"
    
    WORKDIR /app
    
    # Healthcheck
    HEALTHCHECK --interval=30s --timeout=3s CMD nc -zv localhost $PORT || exit 1
    
    # Запуск
    ENTRYPOINT ["/bin/gateway"]
    