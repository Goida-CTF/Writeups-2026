#!/bin/bash
set -e

cd "$(dirname "$0")"
mkdir -p certs
cd certs

# Создаём CA (корневой сертификат)
openssl genrsa -out ca.key 4096
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj '/CN=Ruchkami Dev CA/O=Ruchkami/C=RU'

# Создаём ключ для сервера
openssl genrsa -out server.key 2048

# Создаём конфиг для SAN (Subject Alternative Names)
cat > server.ext << 'EOF'
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = tasks.goidactf.ru
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

# Создаём CSR
openssl req -new -key server.key -out server.csr -subj '/CN=localhost/O=Ruchkami/C=RU'

# Подписываем сертификат нашим CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 825 -sha256 -extfile server.ext

echo "Сертификаты успешно созданы в папке certs/"
ls -la
