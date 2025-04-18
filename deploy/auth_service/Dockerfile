# --- Build stage ---
    FROM golang:1.24.2-alpine AS builder
    WORKDIR /app
    
    # Загрузка зависимостей
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Копируем исходники и билдим
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth_service cmd/auth_service/main.go
    
    # --- Final stage ---
    FROM alpine:3.19
    
    # Устанавливаем зависимости
    RUN apk add --no-cache ca-certificates
    
    # Копируем бинарник
    COPY --from=builder /bin/auth_service /bin/auth_service
    
    # Копируем конфиги
    COPY ./config /app/config
    
    # Runtime переменные
    ENV ENV=dev \
        PORT=50056 \
        LOG_LEVEL=debug \
        AUTH_SECRET_KEY="supersecretkey" \
        TOKEN_TTL=3600
    
    WORKDIR /app
    
    # Healthcheck — например, можно проверить, что порт доступен (если нет отдельного endpoint'а)
    HEALTHCHECK --interval=30s --timeout=3s CMD nc -z localhost $PORT || exit 1
    
    # Запуск
    ENTRYPOINT ["/bin/auth_service"]
    