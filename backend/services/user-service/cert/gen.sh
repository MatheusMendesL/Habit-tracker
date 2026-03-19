rm -f *.pem *.csr *.srl *.cnf

openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
-keyout ca-key.pem \
-out ca-cert.pem

openssl genrsa -out user-service-key.pem 2048

openssl req -new \
-key user-service-key.pem \
-out user-service.csr

cat > user-service-ext.cnf <<EOF
subjectAltName=DNS:localhost,IP:127.0.0.1
EOF

openssl x509 -req \
-in user-service.csr \
-CA ca-cert.pem \
-CAkey ca-key.pem \
-CAcreateserial \
-out user-service-cert.pem \
-days 365 \
-extfile user-service-ext.cnf