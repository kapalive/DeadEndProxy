#!/usr/bin/env sh
# © 2023 Devinsidercode CORP. Licensed under the MIT License.
#
# Helper script for Linux systems to open the configuration
# file in the default editor or run the proxy binary.
CONFIG_PATH="/etc/deadendproxy/config.yaml"
BINARY="/usr/local/bin/deadendproxy-bin"

if [ "$1" = "config" ]; then
    ${EDITOR:-nano} "$CONFIG_PATH"
else
    exec "$BINARY" "$@"
fi
