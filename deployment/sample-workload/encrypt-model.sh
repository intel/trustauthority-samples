#!/bin/bash

set -x

public_key=$(cat public.crt)

# Get DEK from KBS using public key
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

# Check the exit status of the curl command
if [ $? -ne 0 ]; then
    echo "Error: Failed to retrieve token from KBS"
    exit 1
fi

# Create Key transfer policy
KEY_TRANSFER_OUTPUT=$(curl --location ${KBS_URL}/key-transfer-policies \
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
    echo "Error: Failed to create key transfer policy"
    exit 1
fi

# Extract transfer_policy_id from key transfer policy creation output
transfer_policy_id=$(echo "$KEY_TRANSFER_OUTPUT" | jq -r '.id')
if [ $? -ne 0 ]; then
    echo "Error: Failed to retrieve transfer_policy_id from key transfer api output"
    exit 1
fi

# Create AES 256 Key on KBS with transfer_policy_id
KEY_OUTPUT=$(curl --location ${KBS_URL}/keys \
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
    echo "Error: Failed to create key"
    exit 1
fi

# Extract id from key creation output
id=$(echo "$KEY_OUTPUT" | jq -r '.id')
if [ $? -ne 0 ]; then
    echo "Error: Failed to extract id from key creation output"
    exit 1
fi

# Fetch the wrapped key
WRAPPED_KEY_OUTPUT=$(curl --location ${KBS_URL}/keys/$id \
--header 'Content-Type: application/x-pem-file' \
--header 'Accept: application/json' \
--header "Authorization: Bearer ${BEARER_TOKEN}" \
--data "${public_key}" -k)

if [ $? -ne 0 ]; then
    echo "Error: Failed to fetch wrapped DEK"
    exit 1
fi

wrapped_key=$(echo "$WRAPPED_KEY_OUTPUT" | jq -r '.wrapped_key')
if [ $? -ne 0 ]; then
    echo "Error: Failed to extract wrapped DEK"
    exit 1
fi

# Generate a sample model
echo "0.055633 0.154901 0.0167112 0.0281701 0.0399299 0.075899 0.0519244 0.0549757 0.189926" > diabetes-linreg.model

echo "$wrapped_key" > wrappedKey

# Encrypt the datafile using encryptor
/usr/local/bin/encrypt diabetes-linreg.model keypair.pem wrappedKey

# Push the encrypted datafile under /etc/
cp model.enc /etc/

echo "Save the key id: "$id" required for /v1/keys API of workload"