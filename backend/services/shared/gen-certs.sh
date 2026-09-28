#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
SERVICES_DIR="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
CERTS_DIR="${SERVICES_DIR}/certs"
CA_KEY="${CERTS_DIR}/ca-key.pem"
CA_CERT="${CERTS_DIR}/ca-cert.pem"
SERVICES=(user-service social-service habit-service stats-service)

umask 077
mkdir -p "${CERTS_DIR}"

openssl req -x509 -newkey rsa:4096 -nodes -sha256 \
  -keyout "${CA_KEY}" \
  -out "${CA_CERT}" \
  -days 3650 \
  -subj "/CN=Habit Tracker Local CA"

for service in "${SERVICES[@]}"; do
  service_cert_dir="${SERVICES_DIR}/${service}/cert"
  service_key="${service_cert_dir}/${service}-key.pem"
  service_csr="${CERTS_DIR}/${service}.csr"
  service_cert="${service_cert_dir}/${service}-cert.pem"
  service_ext="${CERTS_DIR}/${service}.cnf"

  mkdir -p "${service_cert_dir}"
  cp "${CA_CERT}" "${service_cert_dir}/ca-cert.pem"

  openssl genrsa -out "${service_key}" 2048
  openssl req -new \
    -key "${service_key}" \
    -out "${service_csr}" \
    -subj "/CN=${service}"

  cat > "${service_ext}" <<EOF
basicConstraints=CA:FALSE
keyUsage=digitalSignature,keyEncipherment
extendedKeyUsage=serverAuth
subjectAltName=DNS:${service},DNS:${service}.local,DNS:localhost,IP:127.0.0.1
EOF

  openssl x509 -req -sha256 \
    -in "${service_csr}" \
    -CA "${CA_CERT}" \
    -CAkey "${CA_KEY}" \
    -CAcreateserial \
    -out "${service_cert}" \
    -days 825 \
    -extfile "${service_ext}"

  rm -f "${service_csr}" "${service_ext}"
done

rm -f "${CERTS_DIR}/ca-cert.srl"
printf 'Generated TLS certificates for: %s\n' "${SERVICES[*]}"
