#!/bin/sh

# This script will generate auth.json and store it in $SKOPEO_CONFIG_DIR for use
# with other Skopeo commands.
#
# If you already have used Skopeo locally you do not need to use this script.
# Instead, set SKOPEO_CONFIG_DIR to the directory on your local system which
# contains auth.json (likely ~/.config/containers).

# SKOPEO_CONFIG_DIR is the path to the local config directory for Skopeo
SKOPEO_CONFIG_DIR="${PWD}/.config/containers"

echo 'Please enter your Docker Hub credentials when prompted'
echo 'If you are using 2FA then you must generate and use an access token instead of a password'
echo 'Generate an access token at https://hub.docker.com/settings/security'

mkdir -p "$SKOPEO_CONFIG_DIR"

docker run \
    -v "${SKOPEO_CONFIG_DIR}:/tmp/config" \
    -e REGISTRY_AUTH_FILE="/tmp/config/auth.json" \
    -it quay.io/skopeo/stable:latest login docker.io
