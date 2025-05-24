FROM golang:1.24.2-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/user_database ./cmd/user_database

FROM alpine:3.19
RUN apk add --no-cache ca-certificates postgresql-client bash

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY --from=builder /bin/user_database /bin/user_database

COPY --from=builder /app/config /app/config
COPY --from=builder /app/tls /app/tls

WORKDIR /app

ENV ENV=dev \
    PORT=50052 \
    LOG_LEVEL=debug \
    DATABASE_URL="postgres://postgres:password@postgres:5432/pure?sslmode=disable" \
    CONFIG_PATH=/app/config/user_database/config.dev.yaml

RUN ls -la /app/tls && \
    ls -la /app/config

ENTRYPOINT ["/bin/sh", "-c", "/wait-for-it.sh postgres:5432 -t 60 -- /bin/user_database"]
