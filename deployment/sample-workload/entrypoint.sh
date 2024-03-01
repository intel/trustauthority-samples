#!/bin/bash

# Script to encrypt the model
/app/encrypt-model.sh
if [ $? -ne 0 ]; then
    echo "Error: Failed to encrypt model"
    exit 1
fi

# Generate config.json for tdx cli
cat <<EOF > /app/config.json
{
  "trustauthority_api_url": "$TRUSTAUTHORITY_API_URL",
  "trustauthority_api_key": "$TRUSTAUTHORITY_API_KEY"
}
EOF

# Start the workload
/app/trustauthority-demo