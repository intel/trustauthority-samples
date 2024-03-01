#!/bin/bash

openssl genpkey -algorithm RSA -out ca_key.pem

# Create CA certificate
openssl req -new -x509 -sha256 -key ca_key.pem \
       -out ca_certificate.pem -subj "/CN=pykmip CA/O=ITA Intel/C=US"

# Generate client key
openssl genpkey -algorithm RSA -out client_key.pem

# Create OpenSSL configuration file for the client certificate request
cat <<EOF > openssl.cnf
[req]
distinguished_name = req_distinguished_name
req_extensions = req_ext

[req_distinguished_name]

[req_ext]
subjectAltName = IP:127.0.0.1
extendedKeyUsage = serverAuth, clientAuth, codeSigning, emailProtection
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
EOF

# Create client certificate signing request
openssl req -new -sha256 -key client_key.pem -subj "/CN=pykmip" \
       -reqexts req_ext -config openssl.cnf \
       -out client_csr.pem

# Sign client certificate with CA
openssl x509 -req -in client_csr.pem -CA ca_certificate.pem -CAkey ca_key.pem \
       -CAcreateserial -out client_certificate.pem -extfile openssl.cnf -extensions req_ext

# Remove temporary OpenSSL configuration file
rm openssl.cnf
