FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/jwt_service cmd/jwt_service/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/jwt_service /bin/jwt_service
COPY ./config /app/config
ENV ENV=dev \
    PORT=50055 \
    LOG_LEVEL=debug

WORKDIR /app
HEALTHCHECK --interval=30s --timeout=3s CMD nc -z localhost $PORT || exit 1
ENTRYPOINT ["/bin/jwt_service"]
