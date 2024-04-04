#!/bin/bash

echo "Encrypting the model before use"
# Remove extra "/" at end of url
KBS_URL="${KBS_URL%/}"

# Generate kbs auth token
BEARER_TOKEN=$(curl -s --location ${KBS_URL}/token \
              --header 'Accept: application/jwt' \
              --header 'Content-Type: application/json' \
              --data '{
                  "username": "'"${KBS_ADMIN}"'",
                  "password": "'"${KBS_PASSWORD}"'"
              }' -k)
if [ $? -ne 0 ]; then
    echo "Error: Failed to retrieve token from KBS"
    exit 1
fi

# Create Key transfer policy
KEY_TRANSFER_OUTPUT=$(curl -s --location ${KBS_URL}/key-transfer-policies \
                      --header 'Content-Type: application/json' \
                      --header 'Accept: application/json' \
                      --header "Authorization: Bearer ${BEARER_TOKEN}" \
                      --data '{
                          "attestation_type": "TDX",
                          "tdx": {
                              "attributes": {
                                  "enforce_tcb_upto_date": false
                              }
                          }
                      }' -k)
if [ $? -ne 0 ]; then
    echo "Error: Failed to create key transfer policy on KBS"
    exit 1
fi

# Extract transfer_policy_id from key transfer policy creation output
transfer_policy_id=$(echo "$KEY_TRANSFER_OUTPUT" | jq -r '.id')
if [ $? -ne 0 ]; then
    echo "Error: Failed to extract transfer_policy_id from key transfer policy creation output"
    exit 1
fi

# Create AES 256 Key on KBS with transfer_policy_id
KEY_OUTPUT=$(curl -s --location ${KBS_URL}/keys \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header "Authorization: Bearer ${BEARER_TOKEN}" \
--data '{
    "key_information": {
        "algorithm": "AES",
        "key_length": 256
    },
    "transfer_policy_id": "'"$transfer_policy_id"'"
}' -k )
if [ $? -ne 0 ]; then
    echo "Error: Failed to create key on KBS"
    exit 1
fi

# Extract key_id from key creation output
key_id=$(echo "$KEY_OUTPUT" | jq -r '.id')
if [ $? -ne 0 ]; then
    echo "Error: Failed to extract key_id from key creation output"
    exit 1
fi

# Read public key from file
public_key=$(cat public.crt)

# Fetch the wrapped key
WRAPPED_KEY_OUTPUT=$(curl -s --location ${KBS_URL}/keys/$key_id \
--header 'Content-Type: application/x-pem-file' \
--header 'Accept: application/json' \
--header "Authorization: Bearer ${BEARER_TOKEN}" \
--data "${public_key}" -k)
if [ $? -ne 0 ]; then
    echo "Error: Failed to retrieve DEK from KBS"
    exit 1
fi

# Extract DEK from KBS retrieval output
wrapped_key=$(echo "$WRAPPED_KEY_OUTPUT" | jq -r '.wrapped_key')
if [ $? -ne 0 ]; then
    echo "Error: Failed to extract DEK from key retrieval output"
    exit 1
fi

# Write wrapped dek to file
echo "$wrapped_key" > wrappedKey

# Encrypt the datafile using encryptor
/usr/local/bin/encrypt diabetes-linreg.model keypair.pem wrappedKey

# Push the encrypted datafile under /etc/
cp model.enc /etc/

echo "Note the key_id: "$key_id" needed for model decryption later in /v1/key end-point of workload"
