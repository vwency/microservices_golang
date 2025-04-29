FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/gateway cmd/gateway/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates bash postgresql-client
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh
COPY --from=builder /bin/gateway /bin/gateway
COPY ./config/gateway /app/config/gateway

ENV ENV=dev \
    PORT=50058 \
    LOG_LEVEL=debug \
    AUTH_SERVICE_URL="auth_service:50056" \
    DATABASE_URL="postgres://postgres:password@postgres:5432/pure?sslmode=disable"  

WORKDIR /app
HEALTHCHECK --interval=30s --timeout=3s CMD nc -zv localhost $PORT || exit 1
ENTRYPOINT ["/bin/sh", "-c", "/wait-for-it.sh auth_service:50056 -t 6 -- /wait-for-it.sh postgres:5432 -t 30 -- /bin/gateway"]
