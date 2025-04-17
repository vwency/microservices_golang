# Dockerfile
FROM docker.io/library/golang:1.24.2-alpine

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git make protobuf-dev bash curl protoc

# Устанавливаем Task v3
RUN curl -sL https://taskfile.dev/install.sh | sh -s -- -b /usr/local/bin v3.31.0

# Устанавливаем protoc генераторы
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Добавляем GOPATH в PATH
ENV PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /app
COPY . .

# Даем права на выполнение скриптов
RUN chmod +x *.sh

# Выполняем задачи
RUN task proto-generate && \
    task update-deps

EXPOSE 5000-6000

# Запускаем сервисы через Taskfile
CMD ["sh", "-c", "task run-database & task run-auth & task run-gateway"]