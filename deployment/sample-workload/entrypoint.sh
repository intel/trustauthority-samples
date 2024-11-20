#!/bin/bash

# Script to encrypt the model
/app/encrypt-model.sh
if [ $? -ne 0 ]; then
    echo "Error: Failed to encrypt model"
    exit 1
fi

# Install tdx-cli
if [[ -z "$PLATFORM" ]] ; then
    curl -L https://raw.githubusercontent.com/intel/trustauthority-client-for-go/main/release/install-tdx-cli.sh | CLI_VERSION=v1.6.1 bash -
else
    curl -L https://raw.githubusercontent.com/intel/trustauthority-client-for-go/main/release/install-tdx-cli-azure.sh | CLI_VERSION=v1.6.1 bash -
fi
if [ $? -ne 0 ]; then
    echo "Error: Failed to install client"
    exit 1
fi

# Generate config.json for tdx cli
cat <<EOF > /app/config.json
{
  "cloud_provider": "$PLATFORM",
  "trustauthority_api_url": "$TRUSTAUTHORITY_API_URL",
  "trustauthority_api_key": "$TRUSTAUTHORITY_API_KEY"
}
EOF

# Start the workload
/app/trustauthority-demo