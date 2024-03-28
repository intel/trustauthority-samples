#!/bin/bash

COMPONENT_NAME=trustauthority-demo
SERVICE_USERNAME=trustauthority
SERVICE_ENV=trustauthority-demo.env

if [[ $EUID -ne 0 ]]; then
    echo "This installer must be run as root"
    exit 1
fi

# find .env file
echo PWD IS $(pwd)
if [ -f ~/$SERVICE_ENV ]; then
    echo Reading Installation options from $(realpath ~/$SERVICE_ENV)
    env_file=~/$SERVICE_ENV
elif [ -f ../$SERVICE_ENV ]; then
    echo Reading Installation options from $(realpath ../$SERVICE_ENV)
    env_file=../$SERVICE_ENV
fi

if [ -z $env_file ]; then
    NOSETUP="true"
else
    source $env_file
    env_file_exports=$(cat $env_file | grep -E '^[A-Z0-9_]+\s*=' | cut -d = -f 1)
    if [ -n "$env_file_exports" ]; then eval export $env_file_exports; fi
fi

echo "Setting up Trust Authority Linux User..."
# useradd -M -> this user has no home directory
id -u $SERVICE_USERNAME 2> /dev/null || useradd -M --system --shell /sbin/nologin $SERVICE_USERNAME

echo "Installing Intel Trust Authority Demo Application..."

PRODUCT_HOME=/opt/$COMPONENT_NAME
BIN_PATH=$PRODUCT_HOME/bin

for directory in $PRODUCT_HOME $BIN_PATH; do
    mkdir -p $directory
    if [ $? -ne 0 ]; then
        echo "Cannot create directory: $directory"
        exit 1
    fi
    chown -R $SERVICE_USERNAME:$SERVICE_USERNAME $directory
    chmod 700 $directory
done

# Install systemd script
systemctl stop $COMPONENT_NAME >/dev/null 2>&1
systemctl disable $COMPONENT_NAME.service >/dev/null 2>&1
cp $COMPONENT_NAME.service $TLS_CERT_PATH $TLS_KEY_PATH $PRODUCT_HOME && chown $SERVICE_USERNAME:$SERVICE_USERNAME $PRODUCT_HOME/*

cp $COMPONENT_NAME $BIN_PATH/ && chown $SERVICE_USERNAME:$SERVICE_USERNAME $BIN_PATH/*
chmod 700 $BIN_PATH/*
ln -sfT $BIN_PATH/$COMPONENT_NAME /usr/local/bin/$COMPONENT_NAME

# Enable systemd service
systemctl enable $PRODUCT_HOME/$COMPONENT_NAME.service
systemctl daemon-reload

if [ "${NOSETUP,,}" == "true" ]; then
    echo "No .env file found, skipping startup"
    echo "Update systemd configuration with Trust Authority details for manual startup"
    echo "Installation completed successfully!"
else
    # create config.json
    CONFIG_JSON=$PRODUCT_HOME/config.json
    echo "{\"trustauthority_api_url\": \"$TRUSTAUTHORITY_API_URL\"," > $CONFIG_JSON
    echo " \"trustauthority_api_key\": \"$TRUSTAUTHORITY_API_KEY\"}" >> $CONFIG_JSON
    chown $SERVICE_USERNAME:$SERVICE_USERNAME $CONFIG_JSON
    chmod 600 $CONFIG_JSON

    # update systemd daemon configuration
    SYSD_CONF=$PRODUCT_HOME/$COMPONENT_NAME.service
    echo "Environment=TRUSTAUTHORITY_API_URL=$TRUSTAUTHORITY_API_URL" >> $SYSD_CONF
    echo "Environment=TRUSTAUTHORITY_API_KEY=$TRUSTAUTHORITY_API_KEY" >> $SYSD_CONF
    echo "Environment=SKIP_TLS_VERIFICATION=$SKIP_TLS_VERIFICATION" >> $SYSD_CONF
    echo "Environment=HTTPS_PROXY=$HTTPS_PROXY" >> $SYSD_CONF
    chown $SERVICE_USERNAME:$SERVICE_USERNAME $SYSD_CONF
    systemctl daemon-reload

    systemctl start $COMPONENT_NAME
    echo "Waiting for daemon to settle down before checking status"
    sleep 3
    systemctl status $COMPONENT_NAME 2>&1 > /dev/null
    if [ $? != 0 ]; then
      echo "Installation completed with Errors - $COMPONENT_NAME daemon not started."
      echo "Please check errors in syslog using \`journalctl -u $COMPONENT_NAME\`"
      exit 1
    fi
    echo "$COMPONENT_NAME daemon is running"
    echo "Installation completed successfully!"
fi
