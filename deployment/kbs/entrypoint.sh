#!/bin/bash

# Set log level
export PYKMIP_LOG_LEVEL=${PYKMIP_LOG_LEVEL:-info}

# Start pykmip server running as background process
nohup pykmip-server &

# Wait for pykmip server to start
sleep 2

# Start kbs
/kbs run