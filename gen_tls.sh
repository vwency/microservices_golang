#!/bin/bash
set -e

CERT_DIR="tls"
mkdir -p "$CERT_DIR"

# Генерация корневого CA (Certificate Authority)
echo "Генерация корневого CA сертификата..."
openssl genrsa -out $CERT_DIR/ca.key 4096
openssl req -new -x509 -key $CERT_DIR/ca.key -sha256 -subj "/C=RU/ST=State/L=City/O=Organization/CN=root-ca" -days 3650 -out $CERT_DIR/ca.crt

# 1. Генерация серверного ключа и CSR для Database сервиса
echo "Генерация серверного ключа и запроса на подпись сертификата для Database сервиса..."
openssl genrsa -out $CERT_DIR/db_server.key 4096
openssl req -new -key $CERT_DIR/db_server.key -out $CERT_DIR/db_server.csr -config <(
cat <<-EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
C = RU
ST = State
L = City
O = Organization
CN = database_init_service

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = database_init_service
IP.1 = 127.0.0.1
EOF
)

# Подписание серверного сертификата для Database сервиса
echo "Подписание серверного сертификата для Database сервиса..."
openssl x509 -req -in $CERT_DIR/db_server.csr -CA $CERT_DIR/ca.crt -CAkey $CERT_DIR/ca.key -CAcreateserial \
    -out $CERT_DIR/db_server.crt -days 365 -sha256 -extfile <(
cat <<-EOF
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = database_init_service
IP.1 = 127.0.0.1
EOF
)

# 2. Генерация клиентского ключа и CSR для Auth сервиса (когда он является клиентом для Database)
echo "Генерация клиентского ключа и запроса на подпись сертификата для Auth сервиса как клиента..."
openssl genrsa -out $CERT_DIR/auth_client.key 4096
openssl req -new -key $CERT_DIR/auth_client.key -out $CERT_DIR/auth_client.csr -subj "/C=RU/ST=State/L=City/O=ClientOrg/CN=auth-service-client"

# Подписание клиентского сертификата для Auth сервиса
echo "Подписание клиентского сертификата для Auth сервиса как клиента..."
openssl x509 -req -in $CERT_DIR/auth_client.csr -CA $CERT_DIR/ca.crt -CAkey $CERT_DIR/ca.key -CAcreateserial \
    -out $CERT_DIR/auth_client.crt -days 365 -sha256

# 3. Генерация серверного ключа и CSR для Auth сервиса (когда он является сервером)
echo "Генерация серверного ключа и запроса на подпись сертификата для Auth сервиса как сервера..."
openssl genrsa -out $CERT_DIR/auth_server.key 4096
openssl req -new -key $CERT_DIR/auth_server.key -out $CERT_DIR/auth_server.csr -config <(
cat <<-EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
C = RU
ST = State
L = City
O = Organization
CN = auth-service

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = auth-service
IP.1 = 127.0.0.1
EOF
)

# Подписание серверного сертификата для Auth сервиса
echo "Подписание серверного сертификата для Auth сервиса как сервера..."
openssl x509 -req -in $CERT_DIR/auth_server.csr -CA $CERT_DIR/ca.crt -CAkey $CERT_DIR/ca.key -CAcreateserial \
    -out $CERT_DIR/auth_server.crt -days 365 -sha256 -extfile <(
cat <<-EOF
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = auth-service
IP.1 = 127.0.0.1
EOF
)

# Создаем символические ссылки для универсального именования
echo "Создание символических ссылок для удобства использования..."
ln -sf db_server.key $CERT_DIR/server.key
ln -sf db_server.crt $CERT_DIR/server.crt
ln -sf auth_client.key $CERT_DIR/client.key
ln -sf auth_client.crt $CERT_DIR/client.crt

echo "Сертификаты успешно созданы в директории $CERT_DIR"

# Отображение созданных файлов